package dto

// ArticleInput представляет входные данные для создания или обновления контента.
type ArticleInput struct {
	Title     string `json:"title" binding:"required"` // Заголовок контента
	Text      string `json:"text" binding:"required"`  // Текст контента
	Published bool   `json:"published"`                // Опубликован ли контент
	MediaIDs  []uint `json:"media_ids,omitempty"`      // Список ID медиафайлов (опционально)
	AuthorID  uint   `json:"author_id,omitempty"`
}

// ArticleResponse представляет ответ с данными контента.
type ArticleResponse struct {
	ID        uint         `json:"id"`         // Уникальный идентификатор контента.
	AuthorID  uint         `json:"author_id"`  // Идентификатор автора.
	Title     string       `json:"title"`      // Заголовок контента.
	Text      string       `json:"text"`       // Текст контента.
	Published bool         `json:"published"`  // Опубликован ли контент.
	CreatedAt string       `json:"created_at"` // Дата создания.
	UpdatedAt string       `json:"updated_at"` // Дата обновления.
	Media     []MediaDTO   `json:"media"`      // Прикрепленные медиафайлы.
	Comments  []CommentDTO `json:"comments"`   // Комментарии к контенту.
}

// MediaDTO представляет данные медиафайла.
type MediaDTO struct {
	ID       uint   `json:"id"`        // Уникальный идентификатор медиафайла.
	FilePath string `json:"file_path"` // Путь к файлу.
	FileType string `json:"file_type"` // Тип файла.
	FileSize int64  `json:"file_size"` // Размер файла.
}

// CommentDTO представляет данные комментария.
type CommentDTO struct {
	ID       uint   `json:"id"`        // Уникальный идентификатор комментария.
	ParentID *uint  `json:"parent_id"` // Идентификатор родительского комментария.
	Text     string `json:"text"`      // Текст комментария.
}
