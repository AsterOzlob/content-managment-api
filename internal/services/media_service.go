package services

import (
	"errors"
	"os"

	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	"github.com/AsterOzlob/content_managment_api/internal/database/repositories"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/sirupsen/logrus"
)

// MediaService предоставляет методы для управления медиафайлами.
type MediaService struct {
	repo   *repositories.MediaRepository
	Logger logger.Logger
}

// NewMediaService создаёт новый экземпляр MediaService.
func NewMediaService(repo *repositories.MediaRepository, logger logger.Logger) *MediaService {
	return &MediaService{repo: repo, Logger: logger}
}

// UploadFile загружает медиафайл.
func (s *MediaService) UploadFile(input dto.MediaInput, filePath, fileType string, fileSize int64) (*models.Media, error) {
	s.Logger.WithField("file_path", filePath).Info("Uploading media file")

	media := &models.Media{
		ArticleID: input.ArticleID,
		FilePath:  filePath,
		FileType:  fileType,
		FileSize:  fileSize,
	}

	if err := s.repo.Create(media); err != nil {
		s.Logger.WithError(err).Error("Failed to upload media file")
		return nil, err
	}

	return media, nil
}

// GetAllMedia возвращает список всех медиафайлов.
func (s *MediaService) GetAllMedia() ([]*models.Media, error) {
	s.Logger.Info("Fetching all media in service")

	media, err := s.repo.GetAll()
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch all media from repository")
		return nil, err
	}

	return media, nil
}

// DeleteFile удаляет медиафайл по его ID.
func (s *MediaService) DeleteFile(id uint) error {
	s.Logger.WithField("media_id", id).Info("Deleting media file in service")

	media, err := s.repo.GetByID(id)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch media by ID from repository")
		return errors.New("media not found")
	}

	if err := s.repo.Delete(id); err != nil {
		s.Logger.WithError(err).Error("Failed to delete media from repository")
		return err
	}

	if err := os.Remove(media.FilePath); err != nil {
		s.Logger.WithFields(logrus.Fields{
			"file_path": media.FilePath,
			"error":     err.Error(),
		}).Warn("Failed to delete media file from filesystem")
	}

	return nil
}

// GetAllByArticleID возвращает все медиафайлы, связанные с конкретной статьей.
func (s *MediaService) GetAllByArticleID(articleID uint) ([]*models.Media, error) {
	s.Logger.WithField("article_id", articleID).Info("Fetching media by article ID in service")

	media, err := s.repo.GetAllByArticleID(articleID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch media by article ID from repository")
		return nil, err
	}

	return media, nil
}
