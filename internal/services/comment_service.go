package services

import (
	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/models"
	"github.com/AsterOzlob/content_managment_api/internal/repositories"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/sirupsen/logrus"
)

// CommentService предоставляет методы для управления комментариями.
type CommentService struct {
	repo   *repositories.CommentRepository
	Logger *logging.Logger
}

// NewCommentService создает новый экземпляр CommentService.
func NewCommentService(repo *repositories.CommentRepository, logger *logging.Logger) *CommentService {
	return &CommentService{repo: repo, Logger: logger}
}

// AddCommentToArticle добавляет комментарий к статье.
func (s *CommentService) AddCommentToArticle(articleID uint, input dto.CommentInput) (*models.Comment, error) {
	s.Logger.Log(logrus.InfoLevel, "Adding comment to article", map[string]interface{}{
		"article_id": articleID,
	})

	comment := &models.Comment{
		ParentID:  input.ParentID,
		ArticleID: articleID,
		AuthorID:  input.AuthorID,
		Text:      input.Text,
	}

	if err := s.repo.Create(comment); err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to create comment in repository", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	return comment, nil
}

// GetCommentsByArticleID возвращает все комментарии к статье, включая вложенные.
func (s *CommentService) GetCommentsByArticleID(articleID uint) ([]models.Comment, error) {
	s.Logger.Log(logrus.InfoLevel, "Fetching comments by article ID in service", map[string]interface{}{
		"article_id": articleID,
	})

	comments, err := s.repo.GetByArticleID(articleID)
	if err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to fetch comments by article ID from repository", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	return comments, nil
}

// DeleteComment удаляет комментарий по ID.
func (s *CommentService) DeleteComment(commentID uint) error {
	s.Logger.Log(logrus.InfoLevel, "Deleting comment in service", map[string]interface{}{
		"comment_id": commentID,
	})

	if err := s.repo.Delete(commentID); err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to delete comment from repository", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	return nil
}
