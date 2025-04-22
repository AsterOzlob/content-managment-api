package repositories

import (
	"github.com/AsterOzlob/content_managment_api/internal/models"
	"gorm.io/gorm"
)

// UserRepository предоставляет методы для работы с пользователями в базе данных.
type UserRepository struct {
	DB *gorm.DB // DB - экземпляр подключения к базе данных через GORM.
}

// NewUserRepository создает новый экземпляр UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create создает нового пользователя в базе данных.
// Возвращает ошибку, если сохранение не удалось.
func (r *UserRepository) Create(user *models.User) error {
	result := r.DB.Create(user)
	return result.Error
}

// GetAll возвращает список всех пользователей с предзагруженными ролями.
// Используется Preload("Role") для загрузки связанных данных ролей.
func (r *UserRepository) GetAll() ([]*models.User, error) {
	var users []*models.User
	result := r.DB.Preload("Role").Find(&users)
	return users, result.Error
}

// GetByID возвращает пользователя по его ID с предзагруженной ролью.
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	result := r.DB.Preload("Role").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail возвращает пользователя по его email.
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateRole обновляет роль пользователя.
// Находит пользователя по ID и назначает ему роль по названию роли.
func (r *UserRepository) UpdateRole(userID uint, roleName string) error {
	var user models.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		return err
	}

	var role models.Role
	if err := r.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	user.RoleID = role.ID
	return r.DB.Save(&user).Error
}

// Delete удаляет пользователя из базы данных по его ID.
func (r *UserRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.User{}, id)
	return result.Error
}
