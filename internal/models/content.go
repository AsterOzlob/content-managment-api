package models

import "time"

// ContentType определяет тип контента (статья или новость).
type ContentType string

const (
	Article ContentType = "article" // Статья.
	News    ContentType = "news"    // Новость.
)

// Content представляет контент (статью или новость).
type Content struct {
	ID        uint        `json:"id" gorm:"primaryKey"`                 // Уникальный идентификатор контента.
	AuthorID  uint        `json:"author_id" gorm:"not null;index"`      // Идентификатор автора контента.
	Title     string      `json:"title" gorm:"not null;size:255"`       // Заголовок контента.
	Text      string      `json:"text" gorm:"not null;type:text"`       // Текст контента.
	Type      ContentType `json:"type" gorm:"not null;size:20;index"`   // Тип контента (статья или новость).
	Published bool        `json:"published" gorm:"default:false;index"` // Опубликован ли контент.
	CreatedAt time.Time   `json:"created_at"`                           // Дата создания записи.
	UpdatedAt time.Time   `json:"updated_at"`                           // Дата последнего обновления записи.
	Comments  []Comment   `json:"comments"`                             // Комментарии к контенту.
	Media     []Media     `json:"media"`                                // Медиафайлы, связанные с контентом.
}
