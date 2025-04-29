package routes

import (
	"github.com/gin-gonic/gin"
)

// RegisterMediaRoutes настраивает маршруты для управления медиафайлами.
func RegisterMediaRoutes(r *gin.Engine, deps *Dependencies) {
	// Группа маршрутов для медиафайлов
	mediaGroup := r.Group("/media")
	{
		// Загрузка нового медиафайла
		mediaGroup.POST("/upload", deps.MediaCtrl.UploadFile)

		// Получение всех медиафайлов
		mediaGroup.GET("", deps.MediaCtrl.GetAllMedia)

		// Получение всех медиафайлов для конкретной статьи
		mediaGroup.GET("/:id", deps.MediaCtrl.GetAllByArticleID)

		// Удаление медиафайла по ID
		mediaGroup.DELETE("/:id", deps.MediaCtrl.DeleteFile)
	}
}
