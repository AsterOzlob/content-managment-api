package models

import "time"

// User представляет пользователя системы.
type User struct {
	ID           uint       `json:"id" gorm:"primaryKey"`                    // Уникальный идентификатор пользователя.
	RoleID       uint       `gorm:"not null;index"`                          // ID роли пользователя
	Role         Role       `gorm:"foreignKey:RoleID"`                       // Связь с ролью
	Username     string     `json:"username" gorm:"unique;not null;size:64"` // Уникальное имя пользователя.
	Email        string     `json:"email" gorm:"unique;size:255"`            // Уникальный email пользователя.
	PasswordHash string     `json:"-" gorm:"not null"`                       // Хэшированный пароль (скрыт из JSON).
	CreatedAt    time.Time  `json:"created_at"`                              // Дата создания записи.
	UpdatedAt    time.Time  `json:"updated_at"`                              // Дата последнего обновления записи.
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`       // Дата удаления записи (если применимо).
}
