package repositories

import (
	"github.com/AsterOzlob/content_managment_api/internal/models"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// MediaRepository предоставляет методы для работы с медиафайлами в базе данных.
type MediaRepository struct {
	DB     *gorm.DB        // DB - экземпляр подключения к базе данных через GORM.
	Logger *logging.Logger // Logger - экземпляр логгера для MediaRepository.
}

// NewMediaRepository создает новый экземпляр MediaRepository.
func NewMediaRepository(db *gorm.DB, logger *logging.Logger) *MediaRepository {
	return &MediaRepository{DB: db, Logger: logger}
}

// Create создает новый медиафайл в базе данных.
func (r *MediaRepository) Create(media *models.Media) error {
	r.Logger.Log(logrus.InfoLevel, "Creating media in database", map[string]interface{}{
		"file_path": media.FilePath,
	})
	result := r.DB.Create(media)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to create media in database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}
	return nil
}

// GetAll возвращает список всех медиафайлов.
func (r *MediaRepository) GetAll() ([]*models.Media, error) {
	r.Logger.Log(logrus.InfoLevel, "Fetching all media from database", nil)
	var media []*models.Media
	result := r.DB.Find(&media)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to fetch all media from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}
	return media, nil
}

// GetByID возвращает медиафайл по его ID.
func (r *MediaRepository) GetByID(id uint) (*models.Media, error) {
	r.Logger.Log(logrus.InfoLevel, "Fetching media by ID from database", map[string]interface{}{
		"media_id": id,
	})
	var media models.Media
	result := r.DB.First(&media, id)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to fetch media by ID from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}
	return &media, nil
}

// Delete удаляет медиафайл по его ID.
func (r *MediaRepository) Delete(id uint) error {
	r.Logger.Log(logrus.InfoLevel, "Deleting media from database", map[string]interface{}{
		"media_id": id,
	})
	result := r.DB.Delete(&models.Media{}, id)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to delete media from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}
	return nil
}

// GetAllByArticleID возвращает все медиафайлы, связанные с конкретной статьей.
func (r *MediaRepository) GetAllByArticleID(articleID uint) ([]*models.Media, error) {
	r.Logger.Log(logrus.InfoLevel, "Fetching media by article ID from database", map[string]interface{}{
		"article_id": articleID,
	})
	var media []*models.Media
	result := r.DB.Where("article_id = ?", articleID).Find(&media)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to fetch media by article ID from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}
	return media, nil
}
