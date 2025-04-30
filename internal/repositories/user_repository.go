package repositories

import (
	"github.com/AsterOzlob/content_managment_api/internal/models"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// UserRepository предоставляет методы для работы с пользователями в базе данных.
type UserRepository struct {
	DB     *gorm.DB        // DB - экземпляр подключения к базе данных через GORM.
	Logger *logging.Logger // Logger - экземпляр логгера для UserRepository.
}

// NewUserRepository создает новый экземпляр UserRepository.
func NewUserRepository(db *gorm.DB, logger *logging.Logger) *UserRepository {
	return &UserRepository{DB: db, Logger: logger}
}

func (r *UserRepository) Create(user *models.User) error {
	r.Logger.Log(logrus.InfoLevel, "Creating new user in database", map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	})

	// Создаем пользователя
	result := r.DB.Create(user)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to create user in database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}

	// Предзагружаем роль пользователя
	if err := r.DB.Preload("Role").First(&user, user.ID).Error; err != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to preload role for user after creation", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	return nil
}

// GetAll возвращает список всех пользователей с предзагруженными ролями.
func (r *UserRepository) GetAll() ([]*models.User, error) {
	r.Logger.Log(logrus.InfoLevel, "Fetching all users from database", nil)

	var users []*models.User
	result := r.DB.Preload("Role").Find(&users)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to fetch all users from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}
	return users, nil
}

// GetByID возвращает пользователя по его ID с предзагруженной ролью.
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	r.Logger.Log(logrus.InfoLevel, "Fetching user by ID from database", map[string]interface{}{
		"user_id": id,
	})

	var user models.User
	result := r.DB.Preload("Role").First(&user, id)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to fetch user by ID from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail возвращает пользователя по его email.
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	r.Logger.Log(logrus.InfoLevel, "Fetching user by email from database", map[string]interface{}{
		"email": email,
	})
	var user models.User
	result := r.DB.Preload("Role").Where("email = ?", email).First(&user)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to fetch user by email from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}
	return &user, nil
}

// GetRoleByName возвращает роль по её имени.
func (r *UserRepository) GetRoleByName(roleName string) (*models.Role, error) {
	r.Logger.Log(logrus.InfoLevel, "Fetching role by name from database", map[string]interface{}{
		"role_name": roleName,
	})
	var role models.Role
	result := r.DB.Where("name = ?", roleName).First(&role)
	if result.Error != nil {
		r.Logger.Log(logrus.WarnLevel, "Failed to fetch role by name from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}
	return &role, nil
}

// UpdateRole обновляет роль пользователя.
func (r *UserRepository) UpdateRole(userID uint, roleName string) error {
	r.Logger.Log(logrus.InfoLevel, "Updating user role in database", map[string]interface{}{
		"user_id": userID,
		"role":    roleName,
	})

	// Находим пользователя по ID
	var user models.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to find user for role update", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	// Используем новый метод GetRoleByName для получения роли
	role, err := r.GetRoleByName(roleName)
	if err != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to find role for role update", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	// Обновляем роль пользователя
	user.RoleID = role.ID
	if err := r.DB.Save(&user).Error; err != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to update user role in database", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	return nil
}

// Delete удаляет пользователя из базы данных по его ID.
func (r *UserRepository) Delete(id uint) error {
	r.Logger.Log(logrus.InfoLevel, "Deleting user from database", map[string]interface{}{
		"user_id": id,
	})

	result := r.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to delete user from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}
	return nil
}
