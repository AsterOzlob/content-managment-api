package services

import (
	"errors"
	"os"

	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	"github.com/AsterOzlob/content_managment_api/internal/database/repositories"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	apperrors "github.com/AsterOzlob/content_managment_api/pkg/errors"
	"github.com/AsterOzlob/content_managment_api/pkg/utils"
)

// MediaService предоставляет методы для управления медиафайлами.
type MediaService struct {
	repo        *repositories.MediaRepository
	articleRepo *repositories.ArticleRepository
	Logger      logger.Logger
}

// NewMediaService создаёт новый экземпляр MediaService.
func NewMediaService(
	repo *repositories.MediaRepository,
	articleRepo *repositories.ArticleRepository,
	logger logger.Logger,
) *MediaService {
	return &MediaService{
		repo:        repo,
		articleRepo: articleRepo,
		Logger:      logger,
	}
}

// UploadFile загружает медиафайл.
func (s *MediaService) UploadFile(
	input dto.UploadMediaInput,
	authorID uint,
	userRoles []string,
) (*models.Media, error) {
	if input.ArticleID != nil {
		// Только если указан ArticleID — проверяем существование статьи и права пользователя
		article, err := s.articleRepo.GetByID(*input.ArticleID)
		if err != nil {
			s.Logger.WithError(err).WithField("article_id", *input.ArticleID).Error("Failed to get article by ID")
			return nil, errors.New(apperrors.ErrArticleNotFound)
		}
		if !utils.IsOwner(article.AuthorID, authorID, userRoles) {
			s.Logger.Warn("Access denied: user is not the author of the article")
			return nil, errors.New(apperrors.ErrAccessDenied)
		}
	}
	media := &models.Media{
		ArticleID: input.ArticleID,
		AuthorID:  authorID,
		FilePath:  input.FilePath,
		FileType:  input.FileType,
		FileSize:  input.FileSize,
	}
	if err := s.repo.Create(media); err != nil {
		return nil, err
	}
	return media, nil
}

// GetAllMedia возвращает список всех медиафайлов.
func (s *MediaService) GetAllMedia() ([]*models.Media, error) {
	media, err := s.repo.GetAll()
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch all media from repository")
		return nil, errors.New(apperrors.ErrInternalServerError)
	}
	return media, nil
}

// GetAllByArticleID возвращает все медиафайлы, связанные с конкретной статьей.
func (s *MediaService) GetAllByArticleID(articleID uint) ([]*models.Media, error) {
	media, err := s.repo.GetAllByArticleID(articleID)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch media by article ID from repository")
		return nil, errors.New(apperrors.ErrInternalServerError)
	}
	return media, nil
}

// DeleteFile удаляет медиафайл по его ID.
func (s *MediaService) DeleteFile(id uint, userID uint, userRoles []string) error {
	media, err := s.repo.GetByID(id)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch media by ID from repository")
		return errors.New(apperrors.ErrMediaNotFound)
	}
	// Проверяем права пользователя: автор файла, модератор или админ
	if !utils.IsOwner(media.AuthorID, userID, userRoles) {
		s.Logger.Warn("Access denied: user is not the owner or doesn't have required role")
		return errors.New(apperrors.ErrAccessDenied)
	}
	if err := s.repo.Delete(id); err != nil {
		s.Logger.WithError(err).Error("Failed to delete media from repository")
		return err
	}
	if err := os.Remove(media.FilePath); err != nil {
		s.Logger.Warn("Failed to delete media file from filesystem")
	}
	return nil
}
