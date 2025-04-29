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
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// MediaController предоставляет методы для управления медиафайлами через HTTP API.
type MediaController struct {
	service     *services.MediaService // service - экземпляр MediaService для выполнения бизнес-логики.
	Logger      *logging.Logger        // Logger - экземпляр логгера для MediaController.
	MediaConfig *config.MediaConfig    // MediaConfig - конфигурация для работы с медиафайлами.
}

// NewMediaController создает новый экземпляр MediaController.
func NewMediaController(
	service *services.MediaService,
	logger *logging.Logger,
	mediaConfig *config.MediaConfig,
) *MediaController {
	return &MediaController{
		service:     service,
		Logger:      logger,
		MediaConfig: mediaConfig,
	}
}

// @Summary Upload a media file
// @Description Upload a new media file and associate it with an article.
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Param article_id formData uint true "Article ID"
// @Param file formData file true "Media File"
// @Success 201 {object} dto.MediaResponse
// @Failure 400 {object} map[string]string
// @Router /media/upload [post]
func (c *MediaController) UploadFile(ctx *gin.Context) {
	// Получаем ID статьи из формы
	articleIDStr := ctx.PostForm("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 64)
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Invalid article ID", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	// Получаем файл из формы
	file, err := ctx.FormFile("file")
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to get file from request", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to get file"})
		return
	}

	// Проверяем размер файла
	if file.Size > c.MediaConfig.MaxSize {
		c.Logger.Log(logrus.WarnLevel, "File size exceeds the limit", map[string]interface{}{
			"file_size": file.Size,
			"max_size":  c.MediaConfig.MaxSize,
		})
		ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file size exceeds the limit"})
		return
	}

	// Проверяем тип файла
	fileType := file.Header.Get("Content-Type")
	if len(c.MediaConfig.AllowedTypes) > 0 {
		allowed := false
		for _, allowedType := range c.MediaConfig.AllowedTypes {
			if fileType == allowedType {
				allowed = true
				break
			}
		}
		if !allowed {
			c.Logger.Log(logrus.WarnLevel, "File type is not allowed", map[string]interface{}{
				"file_type": fileType,
			})
			ctx.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "file type is not allowed"})
			return
		}
	}

	// Генерируем уникальное имя файла
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	filePath := filepath.Join(c.MediaConfig.StoragePath, fileName)

	// Сохраняем файл на сервере
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to save file", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Вызываем сервис для создания записи о медиафайле
	media, err := c.service.UploadFile(dto.MediaInput{
		ArticleID: uint(articleID),
	}, filePath, fileType, file.Size)
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to upload media file", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем модель в DTO и отправляем ответ
	ctx.JSON(http.StatusCreated, mappers.MapToMediaResponse(media))
}

// @Summary Get all media files
// @Description Get a list of all uploaded media files.
// @Tags Media
// @Produce json
// @Success 200 {array} dto.MediaResponse
// @Failure 500 {object} map[string]string
// @Router /media [get]
func (c *MediaController) GetAllMedia(ctx *gin.Context) {
	c.Logger.Log(logrus.InfoLevel, "Fetching all media files", nil)

	// Вызываем сервис для получения всех медиафайлов
	media, err := c.service.GetAllMedia()
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to fetch all media files", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем модели в DTO и отправляем ответ
	ctx.JSON(http.StatusOK, mappers.MapToMediaListResponse(media))
}

// @Summary Get all media files by article ID
// @Description Get all media files associated with a specific article.
// @Tags Media
// @Produce json
// @Param id path uint true "Article ID"
// @Success 200 {array} dto.MediaResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /media/{id} [get]
func (c *MediaController) GetAllByArticleID(ctx *gin.Context) {
	// Получаем ID статьи из параметров пути
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Invalid article ID", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	// Вызываем сервис для получения медиафайлов
	media, err := c.service.GetAllByArticleID(uint(id))
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to fetch media by article ID", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем модели в DTO и отправляем ответ
	ctx.JSON(http.StatusOK, mappers.MapToMediaListResponse(media))
}

// @Summary Delete a media file
// @Description Delete a media file by its ID.
// @Tags Media
// @Produce json
// @Param id path uint true "Media ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /media/{id} [delete]
func (c *MediaController) DeleteFile(ctx *gin.Context) {
	// Получаем ID медиафайла из параметров пути
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Invalid media ID", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid media ID"})
		return
	}

	// Вызываем сервис для удаления медиафайла
	if err := c.service.DeleteFile(uint(id)); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to delete media file", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Отправляем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{"message": "media file deleted"})
}
