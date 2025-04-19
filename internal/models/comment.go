package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ParentID  *uint  `gorm:"index"`
	ContentID uint   `gorm:"not null;index"`
	AuthorID  uint   `gorm:"not null;index"`
	Author    User   `gorm:"foreignKey:AuthorID"`
	Text      string `gorm:"not null;type:text"`
}

// CanEdit проверяет право редактирования
func (c *Comment) CanEdit(user User) bool {
	return user.ID == c.AuthorID || user.Can("comments", PermissionWrite)
}

// CanDelete проверяет право удаления
func (c *Comment) CanDelete(user User) bool {
	return user.ID == c.AuthorID || user.Can("comments", PermissionDelete)
}
