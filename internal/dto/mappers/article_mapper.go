package mappers

import (
	"time"

	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
)

// MapToContentResponse преобразует модель Article в DTO ContentResponse.
func MapToArticleResponse(content *models.Article) *dto.ArticleResponse {
	var mediaDTOs []dto.MediaDTO
	for _, media := range content.Media {
		mediaDTOs = append(mediaDTOs, dto.MediaDTO{
			ID:       media.ID,
			FilePath: media.FilePath,
			FileType: media.FileType,
			FileSize: media.FileSize,
		})
	}

	var commentDTOs []dto.CommentDTO
	for _, comment := range content.Comments {
		commentDTOs = append(commentDTOs, dto.CommentDTO{
			ID:       comment.ID,
			ParentID: comment.ParentID,
			Text:     comment.Text,
		})
	}

	return &dto.ArticleResponse{
		ID:        content.ID,
		AuthorID:  content.AuthorID,
		Title:     content.Title,
		Text:      content.Text,
		Published: content.Published,
		CreatedAt: content.CreatedAt.Format(time.RFC3339),
		UpdatedAt: content.UpdatedAt.Format(time.RFC3339),
		Media:     mediaDTOs,
		Comments:  commentDTOs,
	}
}

// MapToArticleListResponse преобразует список статей в список DTO-ответов.
func MapToArticleListResponse(articles []*models.Article) []*dto.ArticleResponse {
	dtoArticles := make([]*dto.ArticleResponse, 0, len(articles))

	for _, article := range articles {
		dtoArticles = append(dtoArticles, MapToArticleResponse(article))
	}

	return dtoArticles
}
