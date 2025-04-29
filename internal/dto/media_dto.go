package dto

// MediaInput представляет входные данные для загрузки медиафайла.
type MediaInput struct {
	ArticleID uint `json:"article_id"`
}

// MediaResponse представляет ответ для медиафайла.
type MediaResponse struct {
	ID        uint   `json:"id"`
	ArticleID uint   `json:"article_id"`
	FilePath  string `json:"file_path"`
	FileType  string `json:"file_type"`
	FileSize  int64  `json:"file_size"`
	CreatedAt string `json:"created_at"`
}
