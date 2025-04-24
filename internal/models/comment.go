package models

import "time"

// Comment представляет комментарий к контенту.
type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`             // Уникальный идентификатор комментария.
	ParentID  *uint     `json:"parent_id" gorm:"index"`           // Идентификатор родительского комментария (если есть).
	ArticleID uint      `json:"article_id" gorm:"not null;index"` // Идентификатор контента, к которому относится комментарий.
	AuthorID  uint      `json:"author_id" gorm:"not null;index"`  // Идентификатор автора комментария.
	Text      string    `json:"text" gorm:"not null;type:text"`   // Текст комментария.
	CreatedAt time.Time `json:"created_at"`                       // Дата создания комментария.
	UpdatedAt time.Time `json:"updated_at"`                       // Дата последнего обновления комментария.
}
