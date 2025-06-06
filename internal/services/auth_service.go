package services

import (
	"errors"
	"time"

	"github.com/AsterOzlob/content_managment_api/config"
	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	"github.com/AsterOzlob/content_managment_api/internal/database/repositories"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	apperrors "github.com/AsterOzlob/content_managment_api/pkg/errors"
	"github.com/AsterOzlob/content_managment_api/pkg/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo         *repositories.UserRepository
	refreshTokenRepo *repositories.RefreshTokenRepository
	Logger           logger.Logger
	JWTConfig        *config.JWTConfig
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	refreshTokenRepo *repositories.RefreshTokenRepository,
	logger logger.Logger,
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
	// Проверяем, существует ли пользователь с таким email
	user, err := s.userRepo.GetByEmail(input.Email)
	if err == nil {
		// Если пользователь уже существует
		return nil, nil, errors.New(apperrors.ErrUserAlreadyExists)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Если произошла ошибка, отличная от "запись не найдена"
		s.Logger.WithError(err).Error("Failed to check if user exists")
		return nil, nil, errors.New(apperrors.ErrInternalServerError)
	}

	// Хэшируем пароль
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, nil, errors.New(apperrors.ErrFailedToCreateUser)
	}

	// Создаем нового пользователя
	user = &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
	}

	// Присваиваем роль "user" по умолчанию
	roleName := "user"
	var role models.Role
	result := s.userRepo.DB.Where("name = ?", roleName).First(&role)
	if result.Error != nil {
		s.Logger.WithError(result.Error).Error("Failed to assign default role")
		return nil, nil, errors.New(apperrors.ErrFailedToAssignRole)
	}
	user.RoleID = role.ID

	// Создаем пользователя в базе данных
	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, errors.New(apperrors.ErrFailedToCreateUser)
	}

	// Генерируем токены
	accessToken, err := utils.GenerateAccessToken(user.ID, role.Name, s.JWTConfig)
	if err != nil {
		return nil, nil, errors.New(apperrors.ErrFailedToGenerateTokens)
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.JWTConfig)
	if err != nil {
		return nil, nil, errors.New(apperrors.ErrFailedToGenerateTokens)
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
		return nil, nil, errors.New(apperrors.ErrInternalServerError)
	}

	return user, &AuthTokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Login аутентифицирует пользователя и создаёт токены.
func (s *AuthService) Login(input dto.AuthInput) (*models.User, *AuthTokens, error) {
	user, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		s.Logger.Warn("User not found during login")
		return nil, nil, errors.New(apperrors.ErrInvalidCredentials)
	}
	if err := utils.CheckPasswordHash(input.Password, user.PasswordHash); err != nil {
		s.Logger.Warn("Invalid password during authentication")
		return nil, nil, errors.New(apperrors.ErrInvalidCredentials)
	}
	// Проверяем наличие активных токенов
	existingToken, _ := s.refreshTokenRepo.GetActiveTokenByUser(user.ID)
	if existingToken != nil && existingToken.ExpiresAt.After(time.Now()) {
		// Возвращаем существующие токены
		accessToken, err := utils.RegenerateAccessToken(existingToken.UserID, user.Role.Name, s.JWTConfig)
		if err != nil {
			return nil, nil, errors.New(apperrors.ErrFailedToGenerateTokens)
		}
		return user, &AuthTokens{
			AccessToken:  accessToken,
			RefreshToken: existingToken.Token,
		}, nil
	}
	// Создаём новые токены, если старых нет или они истекли
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role.Name, s.JWTConfig)
	if err != nil {
		return nil, nil, errors.New(apperrors.ErrFailedToGenerateTokens)
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.JWTConfig)
	if err != nil {
		return nil, nil, errors.New(apperrors.ErrFailedToGenerateTokens)
	}
	rt := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Duration(s.JWTConfig.RefreshTokenTTL) * time.Minute),
		IP:        input.IP,
		UserAgent: input.UserAgent,
	}
	if err := s.refreshTokenRepo.Create(rt); err != nil {
		return nil, nil, errors.New(apperrors.ErrInternalServerError)
	}
	return user, &AuthTokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
