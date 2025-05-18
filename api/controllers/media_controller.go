package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/AsterOzlob/content_managment_api/config"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/dto/mappers"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	apperrors "github.com/AsterOzlob/content_managment_api/pkg/errors"
	"github.com/AsterOzlob/content_managment_api/pkg/utils"
	"github.com/gin-gonic/gin"
)

// MediaController предоставляет методы для управления медиафайлами через HTTP API.
type MediaController struct {
	service     *services.MediaService
	mediaConfig *config.MediaConfig
}

// NewMediaController создаёт новый экземпляр MediaController.
func NewMediaController(
	service *services.MediaService,
	mediaConfig *config.MediaConfig,
) *MediaController {
	return &MediaController{
		service:     service,
		mediaConfig: mediaConfig,
	}
}

// UploadFileWithArticle загружает файл с привязкой к статье.
func (c *MediaController) UploadFileWithArticle(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUserNotAuthenticated})
		return
	}
	userRoles, err := utils.GetUserRolesFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUserRolesNotFound})
		return
	}
	articleIDStr := ctx.PostForm("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrInvalidArticleID})
		return
	}
	c.uploadFileInternal(ctx, uint(articleID), userID, userRoles)
}

// UploadUnlinkedFile загружает файл без привязки к статье.
func (c *MediaController) UploadUnlinkedFile(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUserNotAuthenticated})
		return
	}
	userRoles, err := utils.GetUserRolesFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUserRolesNotFound})
		return
	}
	c.uploadFileInternal(ctx, 0, userID, userRoles)
}

// uploadFileInternal — внутренняя функция загрузки файла.
func (c *MediaController) uploadFileInternal(
	ctx *gin.Context,
	articleID uint,
	authorID uint,
	userRoles []string,
) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrInvalidFile})
		return
	}
	if file.Size > c.mediaConfig.MaxSize {
		ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": apperrors.ErrFileSizeExceeded})
		return
	}
	fileType := file.Header.Get("Content-Type")
	allowed := false
	for _, allowedType := range c.mediaConfig.AllowedTypes {
		if allowedType == fileType {
			allowed = true
			break
		}
	}
	if !allowed {
		ctx.JSON(http.StatusUnsupportedMediaType, gin.H{"error": apperrors.ErrUnsupportedFileType})
		return
	}
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	filePath := filepath.Join(c.mediaConfig.StoragePath, fileName)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrFailedToSaveFile})
		return
	}
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, mappers.MapToMediaResponse(media))
}

// GetAllMedia возвращает список всех медиафайлов.
func (c *MediaController) GetAllMedia(ctx *gin.Context) {
	media, err := c.service.GetAllMedia()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrInternalServerError})
		return
	}
	ctx.JSON(http.StatusOK, mappers.MapToMediaListResponse(media))
}

// GetAllByArticleID возвращает медиафайлы по ID статьи.
func (c *MediaController) GetAllByArticleID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrInvalidArticleID})
		return
	}
	media, err := c.service.GetAllByArticleID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrInternalServerError})
		return
	}
	ctx.JSON(http.StatusOK, mappers.MapToMediaListResponse(media))
}

// DeleteFile удаляет медиафайл по ID.
func (c *MediaController) DeleteFile(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrInvalidMediaID})
		return
	}
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUserNotAuthenticated})
		return
	}
	userRoles, err := utils.GetUserRolesFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUserRolesNotFound})
		return
	}
	if err := c.service.DeleteFile(uint(id), userID, userRoles); err != nil {
		switch err.Error() {
		case apperrors.ErrAccessDenied:
			ctx.JSON(http.StatusForbidden, gin.H{"error": apperrors.ErrAccessDenied})
		case apperrors.ErrMediaNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": apperrors.ErrMediaNotFound})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrInternalServerError})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "media file deleted successfully"})
}
