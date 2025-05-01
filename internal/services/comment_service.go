package services

import (
	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	"github.com/AsterOzlob/content_managment_api/internal/database/repositories"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
)

// CommentService предоставляет методы для управления комментариями.
type CommentService struct {
	repo   *repositories.CommentRepository
	Logger logger.Logger
}

// NewCommentService создает новый экземпляр CommentService.
func NewCommentService(repo *repositories.CommentRepository, logger logger.Logger) *CommentService {
	return &CommentService{repo: repo, Logger: logger}
}

// AddCommentToArticle добавляет комментарий к статье.
func (s *CommentService) AddCommentToArticle(articleID uint, input dto.CommentInput) (*models.Comment, error) {
	s.Logger.WithField("article_id", articleID).Info("Adding comment to article")

	comment := &models.Comment{
		ParentID:  input.ParentID,
		ArticleID: articleID,
		AuthorID:  input.AuthorID,
		Text:      input.Text,
	}

	if err := s.repo.Create(comment); err != nil {
		s.Logger.WithError(err).Error("Failed to create comment in repository")
		return nil, err
	}

	return comment, nil
}

// GetCommentsByArticleID возвращает все комментарии к статье, включая вложенные.
func (s *CommentService) GetCommentsByArticleID(articleID uint) ([]*models.Comment, error) {
	s.Logger.WithField("article_id", articleID).Info("Fetching comments by article ID in service")

	comments, err := s.repo.GetByArticleID(articleID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch comments by article ID from repository")
		return nil, err
	}

	return comments, nil
}

// DeleteComment удаляет комментарий по ID.
func (s *CommentService) DeleteComment(commentID uint) error {
	s.Logger.WithField("comment_id", commentID).Info("Deleting comment in service")

	if err := s.repo.Delete(commentID); err != nil {
		s.Logger.WithError(err).Error("Failed to delete comment from repository")
		return err
	}

	return nil
}
