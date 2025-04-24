package repositories

import (
	"github.com/AsterOzlob/content_managment_api/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ContentRepository предоставляет методы для работы с статьями в базе данных.
type ArticleRepository struct {
	DB     *gorm.DB       // DB - экземпляр подключения к базе данных через GORM.
	logger *logrus.Logger // logger - экземпляр логгера для ArticleRepository.
}

// NewArticleRepository создает новый экземпляр ArticleRepository.
func NewArticleRepository(db *gorm.DB, logger *logrus.Logger) *ArticleRepository {
	return &ArticleRepository{DB: db, logger: logger}
}

// Create создает новую статью в базе данных.
func (r *ArticleRepository) Create(article *models.Article) error {
	r.logger.WithFields(logrus.Fields{
		"author_id": article.AuthorID,
		"title":     article.Title,
	}).Info("Creating article in database")

	result := r.DB.Create(article)
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"error": result.Error.Error(),
		}).Error("Failed to create article in database")
		return result.Error
	}
	return nil
}

// GetAll возвращает список всех статей с медиафайлами и комментариями.
func (r *ArticleRepository) GetAll() ([]*models.Article, error) {
	r.logger.Info("Fetching all articles from database")

	var articles []*models.Article
	result := r.DB.Preload("Media").Preload("Comments").Find(&articles)
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"error": result.Error.Error(),
		}).Error("Failed to fetch all articles from database")
		return nil, result.Error
	}
	return articles, nil
}

// GetByID возвращает статью по ее ID с медиафайлами и комментариями.
func (r *ArticleRepository) GetByID(id uint) (*models.Article, error) {
	r.logger.WithFields(logrus.Fields{
		"article_id": id,
	}).Info("Fetching article by ID from database")

	var article models.Article
	result := r.DB.Preload("Media").Preload("Comments").First(&article, id)
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"error": result.Error.Error(),
		}).Error("Failed to fetch article by ID from database")
		return nil, result.Error
	}
	return &article, nil
}

// Update обновляет существующую статью.
func (r *ArticleRepository) Update(article *models.Article) error {
	r.logger.WithFields(logrus.Fields{
		"article_id": article.ID,
		"title":      article.Title,
	}).Info("Updating article in database")

	result := r.DB.Save(article)
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"error": result.Error.Error(),
		}).Error("Failed to update article in database")
		return result.Error
	}
	return nil
}

// Delete удаляет статью по ее ID.
func (r *ArticleRepository) Delete(id uint) error {
	r.logger.WithFields(logrus.Fields{
		"article_id": id,
	}).Info("Deleting article from database")

	result := r.DB.Delete(&models.Article{}, id)
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"error": result.Error.Error(),
		}).Error("Failed to delete article from database")
		return result.Error
	}
	return nil
}

// AttachMedia прикрепляет медиафайлы к статье.
func (r *ArticleRepository) AttachMedia(articleID uint, mediaIDs []uint) error {
	r.logger.WithFields(logrus.Fields{
		"article_id": articleID,
		"media_ids":  mediaIDs,
	}).Info("Attaching media to article")

	for _, mediaID := range mediaIDs {
		media := models.Media{ID: mediaID, ArticleID: articleID}
		result := r.DB.Model(&media).Update("article_id", articleID)
		if result.Error != nil {
			r.logger.WithFields(logrus.Fields{
				"error": result.Error.Error(),
			}).Error("Failed to attach media to article")
			return result.Error
		}
	}
	return nil
}
