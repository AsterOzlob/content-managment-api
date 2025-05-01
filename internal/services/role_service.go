package services

import (
	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	"github.com/AsterOzlob/content_managment_api/internal/database/repositories"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
)

type RoleService struct {
	repo   *repositories.RoleRepository
	Logger logger.Logger
}

// NewRoleService создаёт новый экземпляр RoleService.
func NewRoleService(repo *repositories.RoleRepository, logger logger.Logger) *RoleService {
	return &RoleService{repo: repo, Logger: logger}
}

// CreateRole создаёт новую роль.
func (s *RoleService) CreateRole(input *dto.RoleCreateDTO) (*models.Role, error) {
	s.Logger.WithField("role_name", input.Name).Info("Creating new role via service")

	role := &models.Role{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := s.repo.Create(role); err != nil {
		s.Logger.WithError(err).Error("Failed to create role via service")
		return nil, err
	}

	return role, nil
}

// GetAllRoles возвращает список всех ролей.
func (s *RoleService) GetAllRoles() ([]*models.Role, error) {
	s.Logger.Info("Fetching all roles via service")

	roles, err := s.repo.GetAllRoles()
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch all roles via service")
		return nil, err
	}

	return roles, nil
}

// GetRoleByID получает роль по ID.
func (s *RoleService) GetRoleByID(id uint) (*models.Role, error) {
	s.Logger.WithField("role_id", id).Info("Fetching role by ID via service")

	role, err := s.repo.GetRoleByID(id)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch role by ID via service")
		return nil, err
	}

	return role, nil
}

// UpdateRole обновляет существующую роль.
func (s *RoleService) UpdateRole(id uint, input *dto.RoleUpdateDTO) (*models.Role, error) {
	s.Logger.WithField("role_id", id).Info("Updating role via service")

	role, err := s.repo.GetRoleByID(id)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch role for update via service")
		return nil, err
	}

	if input.Name != "" {
		role.Name = input.Name
	}
	if input.Description != "" {
		role.Description = input.Description
	}

	if err := s.repo.Update(role); err != nil {
		s.Logger.WithError(err).Error("Failed to update role via service")
		return nil, err
	}

	return role, nil
}

// DeleteRole удаляет роль по ID.
func (s *RoleService) DeleteRole(id uint) error {
	s.Logger.WithField("role_id", id).Info("Deleting role via service")

	if err := s.repo.Delete(id); err != nil {
		s.Logger.WithError(err).Error("Failed to delete role via service")
		return err
	}

	return nil
}
