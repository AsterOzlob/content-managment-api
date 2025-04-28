package repositories

import (
	"github.com/AsterOzlob/content_managment_api/internal/models"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CommentRepository предоставляет методы для работы с комментариями в базе данных.
type CommentRepository struct {
	DB     *gorm.DB
	Logger *logging.Logger
}

// NewCommentRepository создает новый экземпляр CommentRepository.
func NewCommentRepository(db *gorm.DB, logger *logging.Logger) *CommentRepository {
	return &CommentRepository{DB: db, Logger: logger}
}

// Create создает новый комментарий в базе данных.
func (r *CommentRepository) Create(comment *models.Comment) error {
	r.Logger.Log(logrus.InfoLevel, "Creating comment in database", map[string]interface{}{
		"article_id": comment.ArticleID,
	})

	result := r.DB.Create(comment)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to create comment in database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}

	return nil
}

// GetByArticleID возвращает все комментарии к статье, включая вложенные.
func (r *CommentRepository) GetByArticleID(articleID uint) ([]models.Comment, error) {
	r.Logger.Log(logrus.InfoLevel, "Fetching comments by article ID from database", map[string]interface{}{
		"article_id": articleID,
	})

	var comments []models.Comment
	result := r.DB.Preload("Replies").Where("article_id = ? AND parent_id IS NULL", articleID).Find(&comments)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to fetch comments by article ID from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}

	return comments, nil
}

// Delete удаляет комментарий по ID.
func (r *CommentRepository) Delete(id uint) error {
	r.Logger.Log(logrus.InfoLevel, "Deleting comment from database", map[string]interface{}{
		"comment_id": id,
	})

	result := r.DB.Delete(&models.Comment{}, id)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to delete comment from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}

	return nil
}
