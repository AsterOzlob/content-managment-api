package services

import (
	"errors"

	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	"github.com/AsterOzlob/content_managment_api/internal/database/repositories"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	apperrors "github.com/AsterOzlob/content_managment_api/pkg/errors"
	"github.com/AsterOzlob/content_managment_api/pkg/utils"
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
func (s *CommentService) AddCommentToArticle(articleID uint, input dto.CommentInput, userID uint) (*models.Comment, error) {
	comment := &models.Comment{
		ParentID:  input.ParentID,
		ArticleID: articleID,
		AuthorID:  userID,
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
	comments, err := s.repo.GetByArticleID(articleID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch comments by article ID from repository")
		return nil, errors.New(apperrors.ErrArticleNotFound)
	}
	return comments, nil
}

// UpdateComment редактирует содержимое комментария.
func (s *CommentService) UpdateComment(id uint, input dto.CommentInput, userID uint, roles []string) (*models.Comment, error) {
	comment, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New(apperrors.ErrCommentNotFound)
	}
	if !utils.IsOwner(comment.AuthorID, userID, roles) {
		return nil, errors.New(apperrors.ErrAccessDenied)
	}
	comment.Text = input.Text
	if err := s.repo.Update(comment); err != nil {
		s.Logger.WithError(err).Error("Failed to update comment in repository")
		return nil, err
	}
	return comment, nil
}

// DeleteComment удаляет комментарий по ID.
func (s *CommentService) DeleteComment(commentID uint, userID uint, userRoles []string) error {
	comment, err := s.repo.GetByID(commentID)
	if err != nil {
		return errors.New(apperrors.ErrCommentNotFound)
	}
	// Проверяем права через IsOwner
	if !utils.IsOwner(comment.AuthorID, userID, userRoles) {
		return errors.New(apperrors.ErrAccessDenied)
	}
	if err := s.repo.Delete(commentID); err != nil {
		s.Logger.WithError(err).Error("Failed to delete comment from repository")
		return err
	}
	return nil
}
