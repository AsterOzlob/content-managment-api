package mappers

import (
	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/models"
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
