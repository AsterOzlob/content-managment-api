package mappers

import (
	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/models"
)

// MapToMediaResponse преобразует модель Media в DTO MediaResponse.
func MapToMediaResponse(media *models.Media) *dto.MediaResponse {
	return &dto.MediaResponse{
		ID:        media.ID,
		ArticleID: media.ArticleID,
		FilePath:  media.FilePath,
		FileType:  media.FileType,
		FileSize:  media.FileSize,
		CreatedAt: media.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// MapToMediaListResponse преобразует список моделей Media в список DTO MediaResponse.
func MapToMediaListResponse(mediaList []*models.Media) []*dto.MediaResponse {
	var dtoMediaList []*dto.MediaResponse
	for _, media := range mediaList {
		dtoMediaList = append(dtoMediaList, MapToMediaResponse(media))
	}
	return dtoMediaList
}
