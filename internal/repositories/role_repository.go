package repositories

import (
	"github.com/AsterOzlob/content_managment_api/internal/models"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RoleRepository предоставляет методы для работы с ролями в базе данных.
type RoleRepository struct {
	DB     *gorm.DB        // DB - экземпляр подключения к базе данных через GORM.
	Logger *logging.Logger // Logger - экземпляр логгера для RoleRepository.
}

// NewRoleRepository создает новый экземпляр RoleRepository.
func NewRoleRepository(db *gorm.DB, logger *logging.Logger) *RoleRepository {
	return &RoleRepository{DB: db, Logger: logger}
}

// CreateRole создает новую роль.
func (r *RoleRepository) Create(role *models.Role) error {
	r.Logger.Log(logrus.InfoLevel, "Creating new role in database", map[string]interface{}{
		"role_name": role.Name,
	})

	result := r.DB.Create(role)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to create role in database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}
	return nil
}

// GetAllRoles возвращает список всех ролей.
func (r *RoleRepository) GetAllRoles() ([]*models.Role, error) {
	r.Logger.Log(logrus.InfoLevel, "Fetching all roles from database", nil)
	var roles []*models.Role
	result := r.DB.Find(&roles)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to fetch all roles from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}
	return roles, nil
}

// GetRoleByID получает роль по ID.
func (r *RoleRepository) GetRoleByID(id uint) (*models.Role, error) {
	r.Logger.Log(logrus.InfoLevel, "Fetching role by ID from database", map[string]interface{}{
		"role_id": id,
	})

	var role models.Role
	result := r.DB.First(&role, id)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to fetch role by ID from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}
	return &role, nil
}

// UpdateRole обновляет существующую роль.
func (r *RoleRepository) Update(role *models.Role) error {
	r.Logger.Log(logrus.InfoLevel, "Updating role in database", map[string]interface{}{
		"role_id": role.ID,
	})

	result := r.DB.Save(role)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to update role in database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}
	return nil
}

// DeleteRole удаляет роль по ID.
func (r *RoleRepository) Delete(id uint) error {
	r.Logger.Log(logrus.InfoLevel, "Deleting role from database", map[string]interface{}{
		"role_id": id,
	})

	result := r.DB.Delete(&models.Role{}, id)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to delete role from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}
	return nil
}
