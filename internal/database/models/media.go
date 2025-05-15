package models

import (
	"time"

	"github.com/AsterOzlob/content_managment_api/pkg/utils"
	"gorm.io/gorm"
)

// Media представляет медиафайл, связанный с контентом.
type Media struct {
	ID        uint      `json:"id" gorm:"primaryKey"`                           // Уникальный идентификатор медиафайла.
	ArticleID *uint     `json:"article_id,omitempty" gorm:"index,default:null"` // Идентификатор контента, к которому относится файл.
	AuthorID  uint      `json:"author_id" gorm:"not null;index"`                // Идентификатор автора файла.
	FilePath  string    `json:"file_path" gorm:"not null"`                      // Путь к файлу на сервере.
	FileType  string    `json:"file_type" gorm:"not null;size:50"`              // Тип файла (например, "image/jpeg").
	FileSize  int64     `json:"file_size" gorm:"not null"`                      // Размер файла в байтах.
	CreatedAt time.Time `json:"created_at"`                                     // Дата загрузки файла.
	UpdatedAt time.Time `json:"updated_at"`                                     // Дата последнего обновления записи.
}

// BeforeCreate вызывается перед сохранением новой записи.
// Предназначена для санитизации строковых полей и защиты от XSS-атак.
func (m *Media) BeforeCreate(tx *gorm.DB) (err error) {
	m.FilePath = utils.Sanitize(m.FilePath)
	return nil
}

// BeforeUpdate вызывается перед обновлением записи.
// Используется для очистки данных от потенциально опасного HTML/JS.
func (m *Media) BeforeUpdate(tx *gorm.DB) (err error) {
	m.FilePath = utils.Sanitize(m.FilePath)
	return nil
}
