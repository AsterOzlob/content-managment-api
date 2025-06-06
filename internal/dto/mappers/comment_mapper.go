package mappers

import (
	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
)

// MapToCommentResponse преобразует модель Comment в DTO CommentResponse.
func MapToCommentResponse(comment *models.Comment) dto.CommentResponse {
	var replies []dto.CommentResponse
	for _, reply := range comment.Replies {
		replies = append(replies, MapToCommentResponse(&reply))
	}

	return dto.CommentResponse{
		ID:        comment.ID,
		ParentID:  comment.ParentID,
		ArticleID: comment.ArticleID,
		AuthorID:  comment.AuthorID,
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		Replies:   replies,
	}
}

// MapToCommentListResponse преобразует список моделей Comment в список DTO CommentResponse.
func MapToCommentListResponse(comments []*models.Comment) []dto.CommentResponse {
	dtoComments := make([]dto.CommentResponse, 0, len(comments))

	for _, comment := range comments {
		dtoComments = append(dtoComments, MapToCommentResponse(comment))
	}

	return dtoComments
}
