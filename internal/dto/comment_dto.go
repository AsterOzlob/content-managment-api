package dto

import "time"

// CommentInput представляет входные данные для создания комментария.
type CommentInput struct {
	ParentID  *uint  `json:"parent_id"`               // ID родительского комментария (опционально).
	ArticleID uint   `json:"article_id"`              // ID статьи.
	AuthorID  uint   `json:"author_id"`               // ID автора.
	Text      string `json:"text" binding:"required"` // Текст комментария.
}

// CommentResponse представляет ответ с данными комментария.
type CommentResponse struct {
	ID        uint              `json:"id"`                // Уникальный идентификатор комментария.
	ParentID  *uint             `json:"parent_id"`         // ID родительского комментария (если есть).
	ArticleID uint              `json:"article_id"`        // ID статьи.
	AuthorID  uint              `json:"author_id"`         // ID автора.
	Text      string            `json:"text"`              // Текст комментария.
	CreatedAt time.Time         `json:"created_at"`        // Дата создания комментария.
	UpdatedAt time.Time         `json:"updated_at"`        // Дата последнего обновления комментария.
	Replies   []CommentResponse `json:"replies,omitempty"` // Вложенные комментарии.
}
