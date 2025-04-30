package services

import (
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
