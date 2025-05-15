package models

import (
	"time"

	"github.com/AsterOzlob/content_managment_api/pkg/utils"
	"gorm.io/gorm"
)

// Comment представляет комментарий к контенту.
type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`             // Уникальный идентификатор комментария.
	ParentID  *uint     `json:"parent_id" gorm:"index"`           // Идентификатор родительского комментария (если есть).
	ArticleID uint      `json:"article_id" gorm:"not null;index"` // Идентификатор контента, к которому относится комментарий.
	AuthorID  uint      `json:"author_id" gorm:"not null;index"`  // Идентификатор автора комментария.
	Text      string    `json:"text" gorm:"not null;type:text"`   // Текст комментария.
	CreatedAt time.Time `json:"created_at"`                       // Дата создания комментария.
	UpdatedAt time.Time `json:"updated_at"`                       // Дата последнего обновления комментария.

	// Вложенные комментарии (рекурсивная связь)
	Replies []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE;"` // Дочерние комментарии.
}

// BeforeCreate вызывается перед сохранением новой записи.
// Предназначена для санитизации строковых полей и защиты от XSS-атак.
func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.Text = utils.Sanitize(c.Text)
	return nil
}

// BeforeUpdate вызывается перед обновлением записи.
// Используется для очистки данных от потенциально опасного HTML/JS.
func (c *Comment) BeforeUpdate(tx *gorm.DB) (err error) {
	c.Text = utils.Sanitize(c.Text)
	return nil
}
