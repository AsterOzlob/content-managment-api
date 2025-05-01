package models

import "time"

// RefreshToken представляет refresh token для аутентификации.
type RefreshToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`                  // Уникальный идентификатор токена.
	UserID    uint      `json:"user_id" gorm:"not null;index"`         // Идентификатор пользователя, которому принадлежит токен.
	Token     string    `json:"token" gorm:"unique;not null;size:512"` // Сам токен.
	UserAgent string    `json:"user_agent" gorm:"size:255"`            // User-Agent клиента.
	IP        string    `json:"ip" gorm:"size:45"`                     // IP-адрес клиента.
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`      // Время истечения токена.
	Revoked   bool      `json:"revoked" gorm:"default:false;index"`    // Отозван ли токен.
	CreatedAt time.Time `json:"created_at"`                            // Дата создания токена.
	UpdatedAt time.Time `json:"updated_at"`                            // Дата последнего обновления записи.
}
