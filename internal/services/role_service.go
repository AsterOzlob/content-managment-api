package services

import (
	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/models"
	"github.com/AsterOzlob/content_managment_api/internal/repositories"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/sirupsen/logrus"
)

type RoleService struct {
	repo   *repositories.RoleRepository
	Logger *logging.Logger
}

// NewRoleService создает новый экземпляр RoleService.
func NewRoleService(repo *repositories.RoleRepository, logger *logging.Logger) *RoleService {
	return &RoleService{repo: repo, Logger: logger}
}

// CreateRole создает новую роль.
func (s *RoleService) CreateRole(input *dto.RoleCreateDTO) (*models.Role, error) {
	s.Logger.Log(logrus.InfoLevel, "Creating new role via service", map[string]interface{}{
		"role_name": input.Name,
	})
	role := &models.Role{
		Name:        input.Name,
		Description: input.Description,
	}
	if err := s.repo.Create(role); err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to create role via service", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	return role, nil
}

// GetAllRoles возвращает список всех ролей.
func (s *RoleService) GetAllRoles() ([]*models.Role, error) {
	s.Logger.Log(logrus.InfoLevel, "Fetching all roles via service", nil)
	roles, err := s.repo.GetAllRoles()
	if err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to fetch all roles via service", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	return roles, nil
}

// GetRoleByID получает роль по ID.
func (s *RoleService) GetRoleByID(id uint) (*models.Role, error) {
	s.Logger.Log(logrus.InfoLevel, "Fetching role by ID via service", map[string]interface{}{
		"role_id": id,
	})

	role, err := s.repo.GetRoleByID(id)
	if err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to fetch role by ID via service", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	return role, nil
}

// UpdateRole обновляет существующую роль.
func (s *RoleService) UpdateRole(id uint, input *dto.RoleUpdateDTO) (*models.Role, error) {
	s.Logger.Log(logrus.InfoLevel, "Updating role via service", map[string]interface{}{
		"role_id": id,
	})
	role, err := s.repo.GetRoleByID(id)
	if err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to fetch role for update via service", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	// Обновляем только те поля, которые пришли в DTO
	if input.Name != "" {
		role.Name = input.Name
	}
	if input.Description != "" {
		role.Description = input.Description
	}
	if err := s.repo.Update(role); err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to update role via service", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	return role, nil
}

// DeleteRole удаляет роль по ID.
func (s *RoleService) DeleteRole(id uint) error {
	s.Logger.Log(logrus.InfoLevel, "Deleting role via service", map[string]interface{}{
		"role_id": id,
	})

	if err := s.repo.Delete(id); err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to delete role via service", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	return nil
}
