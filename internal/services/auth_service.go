package services

import (
	"errors"
	"time"

	"github.com/AsterOzlob/content_managment_api/config"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/models"
	"github.com/AsterOzlob/content_managment_api/internal/repositories"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/AsterOzlob/content_managment_api/utils"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	userRepo         *repositories.UserRepository
	refreshTokenRepo *repositories.RefreshTokenRepository
	Logger           *logging.Logger
	JWTConfig        *config.JWTConfig
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	refreshTokenRepo *repositories.RefreshTokenRepository,
	logger *logging.Logger,
	jwtConfig *config.JWTConfig,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		Logger:           logger,
		JWTConfig:        jwtConfig,
	}
}

// SignUp регистрирует нового пользователя и создает токены.
func (s *AuthService) SignUp(input dto.AuthInput) (*models.User, *AuthTokens, error) {
	s.Logger.Log(logrus.InfoLevel, "Registering new user in service", map[string]interface{}{
		"username": input.Username,
		"email":    input.Email,
	})

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, nil, err
	}

	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
	}

	// Присваиваем роль "user" по умолчанию
	roleName := "user"
	var role models.Role
	result := s.userRepo.DB.Where("name = ?", roleName).First(&role)
	if result.Error != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to assign default role", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, nil, errors.New("failed to assign default role")
	}
	user.RoleID = role.ID

	// Создаем пользователя
	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, err
	}

	// Генерируем токены
	accessToken, err := utils.GenerateAccessToken(user.ID, role.Name, s.JWTConfig)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.JWTConfig)
	if err != nil {
		return nil, nil, err
	}

	// Сохраняем refresh token в базу данных
	rt := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Duration(s.JWTConfig.RefreshTokenTTL) * time.Minute),
		IP:        input.IP,
		UserAgent: input.UserAgent,
	}
	if err := s.refreshTokenRepo.Create(rt); err != nil {
		return nil, nil, err
	}

	return user, &AuthTokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Login аутентифицирует пользователя и создает токены.
func (s *AuthService) Login(input dto.AuthInput) (*models.User, *AuthTokens, error) {
	s.Logger.Log(logrus.InfoLevel, "Authenticating user in service", map[string]interface{}{
		"email": input.Email,
	})

	user, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		s.Logger.Log(logrus.WarnLevel, "User not found during login", nil)
		return nil, nil, errors.New("invalid credentials")
	}

	if err := utils.CheckPasswordHash(input.Password, user.PasswordHash); err != nil {
		s.Logger.Log(logrus.WarnLevel, "Invalid password during authentication", nil)
		return nil, nil, errors.New("invalid credentials")
	}

	// Проверяем наличие активных токенов
	existingToken, _ := s.refreshTokenRepo.GetActiveTokenByUser(user.ID)
	if existingToken != nil && existingToken.ExpiresAt.After(time.Now()) {
		// Возвращаем существующие токены
		accessToken, err := utils.RegenerateAccessToken(existingToken.UserID, user.Role.Name, s.JWTConfig)
		if err != nil {
			return nil, nil, err
		}
		return user, &AuthTokens{
			AccessToken:  accessToken,
			RefreshToken: existingToken.Token,
		}, nil
	}

	// Создаем новые токены, если старых нет или они истекли
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role.Name, s.JWTConfig)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.JWTConfig)
	if err != nil {
		return nil, nil, err
	}

	rt := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Duration(s.JWTConfig.RefreshTokenTTL) * time.Minute),
		IP:        input.IP,
		UserAgent: input.UserAgent,
	}

	if err := s.refreshTokenRepo.Create(rt); err != nil {
		return nil, nil, err
	}

	return user, &AuthTokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// RefreshToken обновляет access token с использованием refresh token.
func (s *AuthService) RefreshToken(refreshTokenString string) (*AuthTokens, error) {
	s.Logger.Log(logrus.InfoLevel, "Refreshing token in service", map[string]interface{}{
		"refresh_token": refreshTokenString,
	})

	// Находим refresh token в базе данных
	rt, err := s.refreshTokenRepo.GetByToken(refreshTokenString)
	if err != nil || rt == nil {
		s.Logger.Log(logrus.WarnLevel, "Invalid or expired refresh token", nil)
		return nil, errors.New("invalid or expired refresh token")
	}

	// Проверяем срок действия токена
	if rt.ExpiresAt.Before(time.Now()) {
		s.Logger.Log(logrus.WarnLevel, "Expired refresh token", nil)
		return nil, errors.New("expired refresh token")
	}

	// Получаем пользователя
	user, err := s.userRepo.GetByID(rt.UserID)
	if err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to fetch user by ID", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("failed to fetch user")
	}

	// Генерируем новые токены
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role.Name, s.JWTConfig)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.ID, s.JWTConfig)
	if err != nil {
		return nil, err
	}

	// Обновляем refresh token в базе данных
	rt.Token = newRefreshToken
	rt.ExpiresAt = time.Now().Add(time.Duration(s.JWTConfig.RefreshTokenTTL) * time.Minute)
	if err := s.refreshTokenRepo.Update(rt); err != nil {
		return nil, err
	}

	return &AuthTokens{AccessToken: accessToken, RefreshToken: newRefreshToken}, nil
}

// Logout отзывает refresh token.
func (s *AuthService) Logout(refreshTokenString string) error {
	s.Logger.Log(logrus.InfoLevel, "Revoking refresh token", map[string]interface{}{
		"refresh_token": refreshTokenString,
	})

	rt, err := s.refreshTokenRepo.GetByToken(refreshTokenString)
	if err != nil || rt == nil {
		s.Logger.Log(logrus.WarnLevel, "Invalid or not found refresh token during logout", nil)
		return errors.New("invalid or not found refresh token")
	}

	rt.Revoked = true
	if err := s.refreshTokenRepo.Update(rt); err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to revoke refresh token", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	return nil
}
