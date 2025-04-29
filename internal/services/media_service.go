package services

import (
	"errors"
	"os"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/models"
	"github.com/AsterOzlob/content_managment_api/internal/repositories"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/sirupsen/logrus"
)

// MediaService предоставляет методы для управления медиафайлами.
type MediaService struct {
	repo   *repositories.MediaRepository // repo - репозиторий для взаимодействия с базой данных.
	Logger *logging.Logger               // Logger - экземпляр логгера для MediaService.
}

// NewMediaService создает новый экземпляр MediaService.
func NewMediaService(repo *repositories.MediaRepository, logger *logging.Logger) *MediaService {
	return &MediaService{repo: repo, Logger: logger}
}

// UploadFile загружает медиафайл.
func (s *MediaService) UploadFile(input dto.MediaInput, filePath, fileType string, fileSize int64) (*models.Media, error) {
	s.Logger.Log(logrus.InfoLevel, "Uploading media file", map[string]interface{}{
		"file_path": filePath,
	})
	media := &models.Media{
		ArticleID: input.ArticleID,
		FilePath:  filePath,
		FileType:  fileType,
		FileSize:  fileSize,
	}
	if err := s.repo.Create(media); err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to upload media file", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	return media, nil
}

// GetAllMedia возвращает список всех медиафайлов.
func (s *MediaService) GetAllMedia() ([]*models.Media, error) {
	s.Logger.Log(logrus.InfoLevel, "Fetching all media in service", nil)
	media, err := s.repo.GetAll()
	if err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to fetch all media from repository", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	return media, nil
}

// DeleteFile удаляет медиафайл по его ID.
func (s *MediaService) DeleteFile(id uint) error {
	s.Logger.Log(logrus.InfoLevel, "Deleting media file in service", map[string]interface{}{
		"media_id": id,
	})
	// Получаем информацию о файле из базы данных
	media, err := s.repo.GetByID(id)
	if err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to fetch media by ID from repository", map[string]interface{}{
			"error": err.Error(),
		})
		return errors.New("media not found")
	}
	// Удаляем запись из базы данных
	if err := s.repo.Delete(id); err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to delete media from repository", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	// Удаляем файл из файловой системы
	if err := os.Remove(media.FilePath); err != nil {
		s.Logger.Log(logrus.WarnLevel, "Failed to delete media file from filesystem", map[string]interface{}{
			"error":     err.Error(),
			"file_path": media.FilePath,
		})
		// Не возвращаем ошибку, если файл уже удален или не существует
	}
	return nil
}

// GetAllByArticleID возвращает все медиафайлы, связанные с конкретной статьей.
func (s *MediaService) GetAllByArticleID(articleID uint) ([]*models.Media, error) {
	s.Logger.Log(logrus.InfoLevel, "Fetching media by article ID in service", map[string]interface{}{
		"article_id": articleID,
	})
	media, err := s.repo.GetAllByArticleID(articleID)
	if err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to fetch media by article ID from repository", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	return media, nil
}
