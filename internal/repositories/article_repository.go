package repositories

import (
	"github.com/AsterOzlob/content_managment_api/internal/models"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ArticleRepository предоставляет методы для работы с статьями в базе данных.
type ArticleRepository struct {
	DB     *gorm.DB        // DB - экземпляр подключения к базе данных через GORM.
	Logger *logging.Logger // Logger - экземпляр логгера для ArticleRepository.
}

// NewArticleRepository создает новый экземпляр ArticleRepository.
func NewArticleRepository(db *gorm.DB, logger *logging.Logger) *ArticleRepository {
	return &ArticleRepository{DB: db, Logger: logger}
}

// Create создает новую статью в базе данных.
func (r *ArticleRepository) Create(article *models.Article) error {
	r.Logger.Log(logrus.InfoLevel, "Creating article in database", map[string]interface{}{
		"author_id": article.AuthorID,
		"title":     article.Title,
	})

	result := r.DB.Create(article)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to create article in database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}
	return nil
}

// GetAll возвращает список всех статей с медиафайлами и комментариями.
func (r *ArticleRepository) GetAll() ([]*models.Article, error) {
	r.Logger.Log(logrus.InfoLevel, "Fetching all articles from database", nil)

	var articles []*models.Article
	result := r.DB.Preload("Media").Preload("Comments").Find(&articles)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to fetch all articles from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}
	return articles, nil
}

// GetByID возвращает статью по ее ID с медиафайлами и комментариями.
func (r *ArticleRepository) GetByID(id uint) (*models.Article, error) {
	r.Logger.Log(logrus.InfoLevel, "Fetching article by ID from database", map[string]interface{}{
		"article_id": id,
	})

	var article models.Article
	result := r.DB.Preload("Media").Preload("Comments").First(&article, id)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to fetch article by ID from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}
	return &article, nil
}

// Update обновляет существующую статью.
func (r *ArticleRepository) Update(article *models.Article) error {
	r.Logger.Log(logrus.InfoLevel, "Updating article in database", map[string]interface{}{
		"article_id": article.ID,
		"title":      article.Title,
	})

	result := r.DB.Save(article)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to update article in database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}
	return nil
}

// Delete удаляет статью по ее ID.
func (r *ArticleRepository) Delete(id uint) error {
	r.Logger.Log(logrus.InfoLevel, "Deleting article from database", map[string]interface{}{
		"article_id": id,
	})

	result := r.DB.Delete(&models.Article{}, id)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to delete article from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}
	return nil
}

// AttachMedia прикрепляет медиафайлы к статье.
func (r *ArticleRepository) AttachMedia(articleID uint, mediaIDs []uint) error {
	r.Logger.Log(logrus.InfoLevel, "Attaching media to article", map[string]interface{}{
		"article_id": articleID,
		"media_ids":  mediaIDs,
	})

	for _, mediaID := range mediaIDs {
		media := models.Media{ID: mediaID, ArticleID: articleID}
		result := r.DB.Model(&media).Update("article_id", articleID)
		if result.Error != nil {
			r.Logger.Log(logrus.ErrorLevel, "Failed to attach media to article", map[string]interface{}{
				"error": result.Error.Error(),
			})
			return result.Error
		}
	}
	return nil
}
