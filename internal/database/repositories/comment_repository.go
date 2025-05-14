package repositories

import (
	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"gorm.io/gorm"
)

// CommentRepository предоставляет методы для работы с комментариями в базе данных.
type CommentRepository struct {
	DB     *gorm.DB
	Logger logger.Logger
}

// NewCommentRepository создает новый экземпляр CommentRepository.
func NewCommentRepository(db *gorm.DB, logger logger.Logger) *CommentRepository {
	return &CommentRepository{DB: db, Logger: logger}
}

// Create создает новый комментарий в базе данных.
func (r *CommentRepository) Create(comment *models.Comment) error {
	result := r.DB.Create(comment)
	if result.Error != nil {
		r.Logger.WithFields(map[string]interface{}{
			"article_id": comment.ArticleID,
		}).WithError(result.Error).Error("Failed to create comment in database")
		return result.Error
	}
	return nil
}

// GetByArticleID возвращает все комментарии к статье, включая вложенные.
func (r *CommentRepository) GetByArticleID(articleID uint) ([]*models.Comment, error) {
	var comments []*models.Comment
	result := r.DB.Preload("Replies").Where("article_id = ? AND parent_id IS NULL", articleID).Find(&comments)
	if result.Error != nil {
		r.Logger.WithField("article_id", articleID).WithError(result.Error).
			Error("Failed to fetch comments by article ID from database")
		return nil, result.Error
	}
	return comments, nil
}

// GetByID возвращает комментарий по ID.
func (r *CommentRepository) GetByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	result := r.DB.First(&comment, id)
	if result.Error != nil {
		r.Logger.WithField("comment_id", id).WithError(result.Error).
			Error("Failed to fetch comment by ID from database")
		return nil, result.Error
	}
	return &comment, nil
}

// Update редактирует содержимое комментария.
func (r *CommentRepository) Update(comment *models.Comment) error {
	result := r.DB.Save(comment)
	if result.Error != nil {
		r.Logger.WithField("comment_id", comment.ID).WithError(result.Error).
			Error("Failed to update comment in database")
		return result.Error
	}
	return nil
}

// Delete удаляет комментарий по ID.
func (r *CommentRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.Comment{}, id)
	if result.Error != nil {
		r.Logger.WithField("comment_id", id).WithError(result.Error).
			Error("Failed to delete comment from database")
		return result.Error
	}
	return nil
}
