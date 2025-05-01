package repositories

import (
	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"gorm.io/gorm"
)

// RoleRepository предоставляет методы для работы с ролями в базе данных.
type RoleRepository struct {
	DB     *gorm.DB
	Logger logger.Logger
}

// NewRoleRepository создаёт новый экземпляр RoleRepository.
func NewRoleRepository(db *gorm.DB, logger logger.Logger) *RoleRepository {
	return &RoleRepository{DB: db, Logger: logger}
}

// CreateRole создаёт новую роль.
func (r *RoleRepository) Create(role *models.Role) error {
	r.Logger.WithField("name", role.Name).Info("Creating new role in database")
	result := r.DB.Create(role)
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to create role in database")
		return result.Error
	}
	return nil
}

// GetAllRoles возвращает список всех ролей.
func (r *RoleRepository) GetAllRoles() ([]*models.Role, error) {
	r.Logger.Info("Fetching all roles from database")
	var roles []*models.Role
	result := r.DB.Find(&roles)
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to fetch all roles from database")
		return nil, result.Error
	}
	return roles, nil
}

// GetRoleByID получает роль по ID.
func (r *RoleRepository) GetRoleByID(id uint) (*models.Role, error) {
	r.Logger.WithField("role_id", id).Info("Fetching role by ID from database")

	var role models.Role
	result := r.DB.First(&role, id)
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to fetch role by ID from database")
		return nil, result.Error
	}
	return &role, nil
}

// UpdateRole обновляет существующую роль.
func (r *RoleRepository) Update(role *models.Role) error {
	r.Logger.WithField("role_id", role.ID).Info("Updating role in database")
	result := r.DB.Save(role)
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to update role in database")
		return result.Error
	}
	return nil
}

// DeleteRole удаляет роль по ID.
func (r *RoleRepository) Delete(id uint) error {
	r.Logger.WithField("role_id", id).Info("Deleting role from database")
	result := r.DB.Delete(&models.Role{}, id)
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to delete role from database")
		return result.Error
	}
	return nil
}
