package services

import (
	"errors"

	"github.com/AsterOzlob/content_managment_api/internal/models"
	"github.com/AsterOzlob/content_managment_api/internal/repositories"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/sirupsen/logrus"
)

// UserService предоставляет методы для работы с пользователями.
type UserService struct {
	repo   *repositories.UserRepository // repo - репозиторий для взаимодействия с базой данных.
	Logger *logging.Logger              // Logger - экземпляр логгера для UserService.
}

// NewUserService создает новый экземпляр UserService.
func NewUserService(repo *repositories.UserRepository, logger *logging.Logger) *UserService {
	return &UserService{repo: repo, Logger: logger}
}

// SignUp регистрирует нового пользователя.
func (s *UserService) SignUp(username, email, password string) (*models.User, error) {
	s.Logger.Log(logrus.InfoLevel, "Registering new user in service", map[string]interface{}{
		"username": username,
		"email":    email,
	})

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
		s.Logger.Log(logrus.ErrorLevel, "Failed to assign default role", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, errors.New("failed to assign default role")
	}

	user.RoleID = role.ID

	// Сохраняем пользователя в базе данных
	if err := s.repo.Create(user); err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to create user in repository", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	return user, nil
}

// Login выполняет аутентификацию пользователя.
func (s *UserService) Login(email, password string) (*models.User, error) {
	s.Logger.Log(logrus.InfoLevel, "Authenticating user in service", map[string]interface{}{
		"email": email,
	})

	// Находим пользователя по email
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to find user by email", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("user not found")
	}

	// Проверяем пароль
	if user.PasswordHash != password {
		s.Logger.Log(logrus.WarnLevel, "Invalid password during authentication", nil)
		return nil, errors.New("invalid password")
	}

	return user, nil
}

// GetAllUsers возвращает список всех пользователей.
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	s.Logger.Log(logrus.InfoLevel, "Fetching all users in service", nil)

	users, err := s.repo.GetAll()
	if err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to fetch all users from repository", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	return users, nil
}

// GetUserByID возвращает пользователя по ID.
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	s.Logger.Log(logrus.InfoLevel, "Fetching user by ID in service", map[string]interface{}{
		"user_id": id,
	})

	user, err := s.repo.GetByID(id)
	if err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to fetch user by ID from repository", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	return user, nil
}

// DeleteUser удаляет пользователя по ID.
func (s *UserService) DeleteUser(targetUserID uint) error {
	s.Logger.Log(logrus.InfoLevel, "Deleting user in service", map[string]interface{}{
		"user_id": targetUserID,
	})

	if err := s.repo.Delete(targetUserID); err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to delete user from repository", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	return nil
}

// AssignRole назначает роль пользователю.
func (s *UserService) AssignRole(targetUserID uint, roleName string) error {
	s.Logger.Log(logrus.InfoLevel, "Assigning role to user in service", map[string]interface{}{
		"user_id": targetUserID,
		"role":    roleName,
	})

	if err := s.repo.UpdateRole(targetUserID, roleName); err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to assign role to user in repository", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	return nil
}
