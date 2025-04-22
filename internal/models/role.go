package models

import "time"

// Role представляет роль в системе.
type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey"`                   // Уникальный идентификатор роли.
	Name        string    `json:"name" gorm:"unique;not null;size:64"`    // Название роли (например, "admin", "user").
	Description string    `json:"description"`                            // Описание роли.
	Users       []User    `gorm:"foreignKey:RoleID" swaggerignore:"true"` // Пользователи связанные с этой ролью
	CreatedAt   time.Time `json:"created_at"`                             // Дата создания записи.
	UpdatedAt   time.Time `json:"updated_at"`                             // Дата последнего обновления записи.
}
