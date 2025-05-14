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

// @Summary Загрузить медиафайл с привязкой к статье
// @Description Загружает новый медиафайл и связывает его со статьёй по ID.
// @Tags Медиафайлы
// @Accept multipart/form-data
// @Produce json
// @Param article_id formData uint true "ID статьи"
// @Param file formData file true "Медиафайл"
// @Security BearerAuth
// @Success 201 {object} dto.MediaResponse
// @Failure 400 {object} map[string]string
// @Router /media/upload [post]
func (c *MediaController) UploadFileWithArticle(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	userRoles, err := utils.GetUserRolesFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "роли пользователя не найдены"})
		return
	}

	articleIDStr := ctx.PostForm("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID статьи"})
		return
	}

	c.uploadFileInternal(ctx, uint(articleID), userID, userRoles)
}

// @Summary Загрузить медиафайл без привязки к статье
// @Description Загружает медиафайл без связи со статьёй.
// @Tags Медиафайлы
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Медиафайл"
// @Security BearerAuth
// @Success 201 {object} dto.MediaResponse
// @Failure 400 {object} map[string]string
// @Router /media/upload/unlinked [post]
func (c *MediaController) UploadUnlinkedFile(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	userRoles, err := utils.GetUserRolesFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "роли пользователя не найдены"})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "не удалось получить файл"})
		return
	}

	if file.Size > c.mediaConfig.MaxSize {
		ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "размер файла превышает допустимый лимит"})
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
		ctx.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "тип файла не поддерживается"})
		return
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	filePath := filepath.Join(c.mediaConfig.StoragePath, fileName)

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сохранить файл"})
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

// @Summary Получить все медиафайлы
// @Description Возвращает список всех загруженных медиафайлов.
// @Tags Медиафайлы
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.MediaResponse
// @Failure 500 {object} map[string]string
// @Router /media [get]
func (c *MediaController) GetAllMedia(ctx *gin.Context) {
	media, err := c.service.GetAllMedia()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mappers.MapToMediaListResponse(media))
}

// @Summary Получить медиафайлы по ID статьи
// @Description Возвращает все медиафайлы, связанные со статьей по её ID.
// @Tags Медиафайлы
// @Produce json
// @Param id path uint true "ID статьи"
// @Security BearerAuth
// @Success 200 {array} dto.MediaResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /media/{id} [get]
func (c *MediaController) GetAllByArticleID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID статьи"})
		return
	}

	media, err := c.service.GetAllByArticleID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mappers.MapToMediaListResponse(media))
}

// @Summary Удалить медиафайл
// @Description Удаляет медиафайл по его уникальному идентификатору.
// @Tags Медиафайлы
// @Produce json
// @Param id path uint true "ID медиафайла"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /media/{id} [delete]
func (c *MediaController) DeleteFile(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID медиафайла"})
		return
	}

	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	userRoles, err := utils.GetUserRolesFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "роли пользователя не найдены"})
		return
	}

	if err := c.service.DeleteFile(uint(id), userID, userRoles); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "медиафайл успешно удалён"})
}
