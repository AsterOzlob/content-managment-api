package repositories

import (
	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ArticleRepository предоставляет методы для работы со статьями в БД.
type ArticleRepository struct {
	DB     *gorm.DB      // DB - экземпляр подключения к базе данных через GORM.
	Logger logger.Logger // Logger - экземпляр логгера для ArticleRepository.
}

// NewArticleRepository создаёт новый экземпляр ArticleRepository.
func NewArticleRepository(db *gorm.DB, logger logger.Logger) *ArticleRepository {
	return &ArticleRepository{DB: db, Logger: logger}
}

// Create создаёт новую статью в БД.
func (r *ArticleRepository) Create(article *models.Article) error {
	r.Logger.WithFields(logrus.Fields{
		"author_id": article.AuthorID,
		"title":     article.Title,
	}).Info("Creating article in database")

	result := r.DB.Create(article)
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to create article in database")
		return result.Error
	}

	return nil
}

// GetAll возвращает список всех статей.
func (r *ArticleRepository) GetAll() ([]*models.Article, error) {
	r.Logger.Info("Fetching all articles from database")

	var articles []*models.Article
	result := r.DB.Preload("Media").Preload("Comments").Find(&articles)
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to fetch all articles from database")
		return nil, result.Error
	}
	return articles, nil
}

// GetByID возвращает статью по ID.
func (r *ArticleRepository) GetByID(id uint) (*models.Article, error) {
	r.Logger.WithField("article_id", id).Info("Fetching article by ID from database")

	var article models.Article
	result := r.DB.Preload("Media").Preload("Comments").First(&article, id)
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to fetch article by ID from database")
		return nil, result.Error
	}
	return &article, nil
}

// Update обновляет статью в БД.
func (r *ArticleRepository) Update(article *models.Article) error {
	r.Logger.WithFields(logrus.Fields{
		"article_id": article.ID,
		"title":      article.Title,
	}).Info("Updating article in database")

	result := r.DB.Save(article)
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to update article in database")
		return result.Error
	}
	return nil
}

// Delete удаляет статью из БД.
func (r *ArticleRepository) Delete(id uint) error {
	r.Logger.WithField("article_id", id).Info("Deleting article from database")

	result := r.DB.Delete(&models.Article{}, id)
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to delete article from database")
		return result.Error
	}
	return nil
}
