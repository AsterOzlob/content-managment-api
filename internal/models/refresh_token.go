package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	Token     string `gorm:"unique;not null;size:512"`
	UserID    uint   `gorm:"not null;index"`
	UserAgent string `gorm:"size:255"`
	IP        string `gorm:"size:45"`
	ExpiresAt uint64 `gorm:"not null;type:Timestamp;index"`
	Revoked   bool   `gorm:"default:false"`
}

// Метод для проверки валидности токена
func (rt *RefreshToken) IsValid() bool {
	return !rt.Revoked && rt.ExpiresAt > uint64(time.Now().UnixNano())
}
