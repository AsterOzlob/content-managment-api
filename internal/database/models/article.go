package models

import (
	"time"

	"github.com/AsterOzlob/content_managment_api/pkg/utils"
	"gorm.io/gorm"
)

// Article представляет контент (статью или новость).
type Article struct {
	ID        uint      `json:"id" gorm:"primaryKey"`                 // Уникальный идентификатор контента.
	AuthorID  uint      `json:"author_id" gorm:"not null;index"`      // Идентификатор автора контента.
	Title     string    `json:"title" gorm:"not null;size:255"`       // Заголовок контента.
	Text      string    `json:"text" gorm:"not null;type:text"`       // Текст контента.
	Published bool      `json:"published" gorm:"default:false;index"` // Опубликован ли контент.
	CreatedAt time.Time `json:"created_at"`                           // Дата создания записи.
	UpdatedAt time.Time `json:"updated_at"`                           // Дата последнего обновления записи.
	Comments  []Comment `json:"comments"`                             // Комментарии к контенту.
	Media     []Media   `json:"media"`                                // Медиафайлы, связанные с контентом.
}

// BeforeCreate вызывается перед сохранением новой записи.
// Предназначена для санитизации строковых полей и защиты от XSS-атак.
func (a *Article) BeforeCreate(tx *gorm.DB) (err error) {
	a.Title = utils.Sanitize(a.Title)
	a.Text = utils.Sanitize(a.Text)
	return nil
}

// BeforeUpdate вызывается перед обновлением записи.
// Используется для очистки данных от потенциально опасного HTML/JS.
func (a *Article) BeforeUpdate(tx *gorm.DB) (err error) {
	a.Title = utils.Sanitize(a.Title)
	a.Text = utils.Sanitize(a.Text)
	return nil
}
