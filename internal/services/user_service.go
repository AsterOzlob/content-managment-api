package services

import (
	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	"github.com/AsterOzlob/content_managment_api/internal/database/repositories"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
)

// UserService предоставляет методы для работы с пользователями.
type UserService struct {
	repo   *repositories.UserRepository
	Logger logger.Logger
}

// NewUserService создаёт новый экземпляр UserService.
func NewUserService(repo *repositories.UserRepository, logger logger.Logger) *UserService {
	return &UserService{repo: repo, Logger: logger}
}

// GetAllUsers возвращает список всех пользователей.
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch all users from repository")
		return nil, err
	}
	return users, nil
}

// GetUserByID возвращает пользователя по ID.
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		s.Logger.WithError(err).WithField("user_id", id).Error("Failed to fetch user by ID from repository")
		return nil, err
	}
	return user, nil
}

// DeleteUser удаляет пользователя по ID.
func (s *UserService) DeleteUser(targetUserID uint) error {
	if err := s.repo.Delete(targetUserID); err != nil {
		s.Logger.WithError(err).WithField("user_id", targetUserID).Error("Failed to delete user from repository")
		return err
	}
	return nil
}

// AssignRole назначает роль пользователю.
func (s *UserService) AssignRole(targetUserID uint, roleName string) error {
	if err := s.repo.UpdateRole(targetUserID, roleName); err != nil {
		s.Logger.WithError(err).
			WithFields(map[string]interface{}{
				"user_id": targetUserID,
				"role":    roleName,
			}).Error("Failed to assign role to user in repository")
		return err
	}
	return nil
}
