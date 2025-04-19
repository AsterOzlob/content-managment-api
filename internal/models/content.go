package models

import "gorm.io/gorm"

type ContentType string

const (
	Article ContentType = "article"
	News    ContentType = "news"
)

type Content struct {
	gorm.Model
	AuthorID  uint        `gorm:"not null; index"`
	Author    User        `gorm:"foreignKey:AuthorID"`
	Title     string      `gorm:"not null;size:255"`
	Text      string      `gorm:"not null;type:text"`
	Type      ContentType `gorm:"not null;size:20; index"`
	Published bool        `gorm:"default:false; index"`
	Comments  []Comment
	Media     []Media
}

// CanEdit проверяет право редактирования
func (c *Content) CanEdit(user User) bool {
	return user.ID == c.AuthorID || user.Can("content", PermissionWrite)
}

// CanDelete проверяет право удаления
func (c *Content) CanDelete(user User) bool {
	return user.ID == c.AuthorID || user.Can("content", PermissionDelete)
}
