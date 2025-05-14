package repositories

import (
	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"gorm.io/gorm"
)

// MediaRepository предоставляет методы для работы с медиафайлами в базе данных.
type MediaRepository struct {
	DB     *gorm.DB
	Logger logger.Logger
}

// NewMediaRepository создает новый экземпляр MediaRepository.
func NewMediaRepository(db *gorm.DB, logger logger.Logger) *MediaRepository {
	return &MediaRepository{DB: db, Logger: logger}
}

// Create создает новый медиафайл в базе данных.
func (r *MediaRepository) Create(media *models.Media) error {
	result := r.DB.Create(media)
	if result.Error != nil {
		r.Logger.WithField("file_path", media.FilePath).WithError(result.Error).
			Error("Failed to create media in database")
		return result.Error
	}
	return nil
}

// GetAll возвращает список всех медиафайлов.
func (r *MediaRepository) GetAll() ([]*models.Media, error) {
	var media []*models.Media
	result := r.DB.Find(&media)
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to fetch all media from database")
		return nil, result.Error
	}
	return media, nil
}

// GetByID возвращает медиафайл по его ID.
func (r *MediaRepository) GetByID(id uint) (*models.Media, error) {
	var media models.Media
	result := r.DB.First(&media, id)
	if result.Error != nil {
		r.Logger.WithField("media_id", id).WithError(result.Error).
			Error("Failed to fetch media by ID from database")
		return nil, result.Error
	}
	return &media, nil
}

// Delete удаляет медиафайл по его ID.
func (r *MediaRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.Media{}, id)
	if result.Error != nil {
		r.Logger.WithField("media_id", id).WithError(result.Error).
			Error("Failed to delete media from database")
		return result.Error
	}
	return nil
}

// GetAllByArticleID возвращает все медиафайлы, связанные с конкретной статьей.
func (r *MediaRepository) GetAllByArticleID(articleID uint) ([]*models.Media, error) {
	var media []*models.Media
	result := r.DB.Where("article_id = ?", articleID).Find(&media)
	if result.Error != nil {
		r.Logger.WithField("article_id", articleID).WithError(result.Error).
			Error("Failed to fetch media by article ID from database")
		return nil, result.Error
	}
	return media, nil
}
