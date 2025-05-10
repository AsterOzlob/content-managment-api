package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/AsterOzlob/content_managment_api/config"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/dto/mappers"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	"github.com/AsterOzlob/content_managment_api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// MediaController предоставляет методы для управления медиафайлами через HTTP API.
type MediaController struct {
	service     *services.MediaService // service - экземпляр MediaService для выполнения бизнес-логики.
	Logger      logger.Logger          // Logger - интерфейсный логгер
	MediaConfig *config.MediaConfig    // MediaConfig - конфигурация для работы с медиафайлами.
}

// NewMediaController создаёт новый экземпляр MediaController.
func NewMediaController(
	service *services.MediaService,
	logger logger.Logger,
	mediaConfig *config.MediaConfig,
) *MediaController {
	return &MediaController{
		service:     service,
		Logger:      logger,
		MediaConfig: mediaConfig,
	}
}

// @Summary Upload a media file with article ID
// @Description Upload a new media file and associate it with an article.
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Param article_id formData uint true "Article ID"
// @Param file formData file true "Media File"
// @Security BearerAuth
// @Success 201 {object} dto.MediaResponse
// @Failure 400 {object} map[string]string
// @Router /media/upload [post]
func (c *MediaController) UploadFileWithArticle(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to get user ID from context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userRoles, err := utils.GetUserRolesFromContext(ctx)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to get user roles from context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	articleIDStr := ctx.PostForm("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 64)
	if err != nil {
		c.Logger.WithError(err).Error("Invalid article ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	c.uploadFileInternal(ctx, uint(articleID), userID, userRoles)
}

// @Summary Upload an unlinked media file
// @Description Upload a new media file without associating it with any article.
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Media File"
// @Security BearerAuth
// @Success 201 {object} dto.MediaResponse
// @Failure 400 {object} map[string]string
// @Router /media/upload/unlinked [post]
func (c *MediaController) UploadUnlinkedFile(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to get user ID from context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userRoles, err := utils.GetUserRolesFromContext(ctx)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to get user roles from context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.uploadFileInternal(ctx, 0, userID, userRoles)
}

// Внутренняя универсальная функция загрузки
func (c *MediaController) uploadFileInternal(
	ctx *gin.Context,
	articleID uint,
	authorID uint,
	userRoles []string,
) {
	file, err := ctx.FormFile("file")
	if err != nil {
		c.Logger.WithError(err).Error("Failed to get file from request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to get file"})
		return
	}

	if file.Size > c.MediaConfig.MaxSize {
		c.Logger.WithFields(logrus.Fields{
			"file_size": file.Size,
			"max_size":  c.MediaConfig.MaxSize,
		}).Warn("File size exceeds the limit")
		ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file size exceeds the limit"})
		return
	}

	fileType := file.Header.Get("Content-Type")
	if len(c.MediaConfig.AllowedTypes) > 0 && !slices.Contains(c.MediaConfig.AllowedTypes, fileType) {
		c.Logger.WithField("file_type", fileType).Warn("File type is not allowed")
		ctx.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "file type is not allowed"})
		return
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	filePath := filepath.Join(c.MediaConfig.StoragePath, fileName)

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		c.Logger.WithError(err).Error("Failed to save uploaded file")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Формируем новый DTO
	var articleIDPtr *uint
	if articleID != 0 {
		articleIDPtr = &articleID
	}

	uploadInput := dto.UploadMediaInput{
		ArticleID: articleIDPtr,
		FilePath:  filePath,
		FileType:  fileType,
		FileSize:  file.Size,
	}

	media, err := c.service.UploadFile(uploadInput, authorID, userRoles)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to upload media file")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, mappers.MapToMediaResponse(media))
}

// @Summary Get all media files
// @Description Get a list of all uploaded media files.
// @Tags Media
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.MediaResponse
// @Failure 500 {object} map[string]string
// @Router /media [get]
func (c *MediaController) GetAllMedia(ctx *gin.Context) {
	c.Logger.Info("Fetching all media files in controller")

	media, err := c.service.GetAllMedia()
	if err != nil {
		c.Logger.WithError(err).Error("Failed to fetch all media files")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, mappers.MapToMediaListResponse(media))
}

// @Summary Get all media files by article ID
// @Description Get all media files associated with a specific article.
// @Tags Media
// @Produce json
// @Param id path uint true "Article ID"
// @Security BearerAuth
// @Success 200 {array} dto.MediaResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /media/{id} [get]
func (c *MediaController) GetAllByArticleID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Logger.WithError(err).Error("Invalid article ID in GetAllByArticleID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	c.Logger.WithField("article_id", id).Info("Fetching media by article ID")

	media, err := c.service.GetAllByArticleID(uint(id))
	if err != nil {
		c.Logger.WithError(err).Error("Failed to fetch media by article ID")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, mappers.MapToMediaListResponse(media))
}

// @Summary Delete a media file
// @Description Delete a media file by its ID.
// @Tags Media
// @Produce json
// @Param id path uint true "Media ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string "Successfully deleted"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /media/{id} [delete]
func (c *MediaController) DeleteFile(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Logger.WithError(err).Error("Invalid media ID in DeleteFile")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid media ID"})
		return
	}

	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to get user ID from context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userRoles, err := utils.GetUserRolesFromContext(ctx)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to get user roles from context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.Logger.WithField("media_id", id).Info("Deleting media file")

	if err := c.service.DeleteFile(uint(id), userID, userRoles); err != nil {
		c.Logger.WithError(err).Error("Failed to delete media file")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "media file deleted"})
}
