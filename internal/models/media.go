package models

import "gorm.io/gorm"

type Media struct {
	gorm.Model
	ContentID uint   `gorm:"not null;index"`
	AuthorID  uint   `gorm:"not null;index"`
	FilePath  string `gorm:"not null"`
	FileType  string `gorm:"not null;size:50"`
	FileSize  int64  `gorm:"not null"`
}

// CanDelete проверяет право удаления
func (m *Media) CanDelete(user User) bool {
	return user.ID == m.AuthorID || user.Can("media", PermissionDelete)
}
