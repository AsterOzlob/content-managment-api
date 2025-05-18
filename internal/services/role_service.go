package services

import (
	"errors"

	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	"github.com/AsterOzlob/content_managment_api/internal/database/repositories"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	apperrors "github.com/AsterOzlob/content_managment_api/pkg/errors"
)

// RoleService предоставляет методы для управления ролями.
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
	role := &models.Role{
		Name:        input.Name,
		Description: input.Description,
	}
	if err := s.repo.Create(role); err != nil {
		s.Logger.WithError(err).WithField("role_name", input.Name).Error("Failed to create role via service")
		return nil, err
	}
	return role, nil
}

// GetAllRoles возвращает список всех ролей.
func (s *RoleService) GetAllRoles() ([]*models.Role, error) {
	roles, err := s.repo.GetAllRoles()
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch all roles via service")
		return nil, errors.New(apperrors.ErrInternalServerError)
	}
	return roles, nil
}

// GetRoleByID получает роль по ID.
func (s *RoleService) GetRoleByID(id uint) (*models.Role, error) {
	role, err := s.repo.GetRoleByID(id)
	if err != nil {
		s.Logger.WithError(err).WithField("role_id", id).Error("Failed to fetch role by ID via service")
		return nil, errors.New(apperrors.ErrRoleNotFound)
	}
	return role, nil
}

// UpdateRole обновляет существующую роль.
func (s *RoleService) UpdateRole(id uint, input *dto.RoleUpdateDTO) (*models.Role, error) {
	role, err := s.repo.GetRoleByID(id)
	if err != nil {
		s.Logger.WithError(err).WithField("role_id", id).Error("Failed to fetch role for update via service")
		return nil, errors.New(apperrors.ErrRoleNotFound)
	}
	if input.Name != "" {
		role.Name = input.Name
	}
	if input.Description != "" {
		role.Description = input.Description
	}
	if err := s.repo.Update(role); err != nil {
		s.Logger.WithError(err).WithField("role_id", id).Error("Failed to update role via service")
		return nil, err
	}
	return role, nil
}

// DeleteRole удаляет роль по ID.
func (s *RoleService) DeleteRole(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		s.Logger.WithError(err).WithField("role_id", id).Error("Failed to delete role via service")
		return err
	}
	return nil
}
