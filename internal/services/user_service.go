package services

import (
	"errors"
	"time"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/models"
	"github.com/AsterOzlob/content_managment_api/internal/repositories"
)

// UserService предоставляет методы для работы с пользователями.
type UserService struct {
	repo *repositories.UserRepository // repo - репозиторий для взаимодействия с базой данных.
}

// NewUserService создает новый экземпляр UserService.
func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// SignUp регистрирует нового пользователя.
func (s *UserService) SignUp(username, email, password string) (*dto.UserResponse, error) {
	// Создаем нового пользователя
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: password,
	}

	// Присваиваем роль "user" по умолчанию
	roleName := "user"
	var role models.Role
	result := s.repo.DB.Where("name = ?", roleName).First(&role)
	if result.Error != nil {
		return nil, errors.New("failed to assign default role")
	}
	user.RoleID = role.ID

	// Сохраняем пользователя в базе данных
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// Возвращаем ответ с данными пользователя
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      role.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// Login выполняет аутентификацию пользователя.
func (s *UserService) Login(email, password string) (*dto.UserResponse, error) {
	// Находим пользователя по email
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Проверяем пароль
	if user.PasswordHash != password {
		return nil, errors.New("invalid password")
	}

	// Возвращаем данные пользователя
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// DeleteUser удаляет пользователя по ID.
func (s *UserService) DeleteUser(targetUserID uint) error {
	// Удаление пользователя без проверки ролей
	return s.repo.Delete(targetUserID)
}

// GetAllUsers возвращает список всех пользователей.
func (s *UserService) GetAllUsers() ([]*dto.UserResponse, error) {
	// Получаем всех пользователей из базы данных
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// Преобразуем список пользователей в формат DTO
	var result []*dto.UserResponse
	for _, user := range users {
		result = append(result, &dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role.Name,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		})
	}
	return result, nil
}

// GetUserByID возвращает пользователя по ID.
func (s *UserService) GetUserByID(id uint) (*dto.UserResponse, error) {
	// Находим пользователя по ID
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Возвращаем данные пользователя
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// AssignRole назначает роль пользователю.
func (s *UserService) AssignRole(targetUserID uint, roleName string) error {
	// Назначение роли без проверки прав (роли назначаются администратором)
	return s.repo.UpdateRole(targetUserID, roleName)
}
