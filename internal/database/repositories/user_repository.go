package repositories

import (
	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"gorm.io/gorm"
)

// UserRepository предоставляет методы для работы с пользователями в базе данных.
type UserRepository struct {
	DB     *gorm.DB
	Logger logger.Logger
}

// NewUserRepository создаёт новый экземпляр UserRepository.
func NewUserRepository(db *gorm.DB, logger logger.Logger) *UserRepository {
	return &UserRepository{DB: db, Logger: logger}
}

// Create создаёт нового пользователя в базе данных.
func (r *UserRepository) Create(user *models.User) error {
	result := r.DB.Create(user)
	if result.Error != nil {
		r.Logger.WithFields(map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
		}).WithError(result.Error).Error("Failed to create user in database")
		return result.Error
	}
	// Предзагружаем роль пользователя
	if err := r.DB.Preload("Role").First(&user, user.ID).Error; err != nil {
		r.Logger.WithError(err).Error("Failed to preload role for user after creation")
		return err
	}
	return nil
}

// GetAll возвращает список всех пользователей с предзагруженными ролями.
func (r *UserRepository) GetAll() ([]*models.User, error) {
	var users []*models.User
	result := r.DB.Preload("Role").Find(&users)
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to fetch all users from database")
		return nil, result.Error
	}
	return users, nil
}

// GetByID возвращает пользователя по его ID с предзагруженной ролью.
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	result := r.DB.Preload("Role").First(&user, id)
	if result.Error != nil {
		r.Logger.WithField("user_id", id).WithError(result.Error).Error("Failed to fetch user by ID from database")
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail возвращает пользователя по его email.
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.DB.Preload("Role").Where("email = ?", email).First(&user)
	if result.Error != nil {
		r.Logger.WithField("email", email).WithError(result.Error).Error("Failed to fetch user by email from database")
		return nil, result.Error
	}
	return &user, nil
}

// GetRoleByName возвращает роль по её имени.
func (r *UserRepository) GetRoleByName(roleName string) (*models.Role, error) {
	var role models.Role
	result := r.DB.Where("name = ?", roleName).First(&role)
	if result.Error != nil {
		r.Logger.WithField("role_name", roleName).WithError(result.Error).Warn("Failed to fetch role by name from database")
		return nil, result.Error
	}
	return &role, nil
}

// UpdateRole обновляет роль пользователя.
func (r *UserRepository) UpdateRole(userID uint, roleName string) error {
	// Находим пользователя по ID
	var user models.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		r.Logger.WithField("user_id", userID).WithError(err).Error("Failed to find user for role update")
		return err
	}
	// Используем новый метод GetRoleByName для получения роли
	role, err := r.GetRoleByName(roleName)
	if err != nil {
		r.Logger.WithField("role_name", roleName).WithError(err).Error("Failed to find role for role update")
		return err
	}
	// Обновляем роль пользователя
	user.RoleID = role.ID
	if err := r.DB.Save(&user).Error; err != nil {
		r.Logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"role":    roleName,
		}).WithError(err).Error("Failed to update user role in database")
		return err
	}
	return nil
}

// Delete удаляет пользователя из БД.
func (r *UserRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		r.Logger.WithField("user_id", id).WithError(result.Error).Error("Failed to delete user from database")
		return result.Error
	}
	return nil
}
