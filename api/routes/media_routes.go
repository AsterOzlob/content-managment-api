package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterMediaRoutes настраивает маршруты для управления медиафайлами.
func RegisterMediaRoutes(r *gin.Engine, deps *Dependencies) {
	mediaGroup := r.Group("/media")
	{
		// Защищенные эндпоинты
		protected := mediaGroup.Group("/")
		protected.Use(middleware.AuthMiddleware(deps.JWTConfig)) // Middleware для JWT-аутентификации
		{
			// Авторы и администраторы могут загружать файлы
			protected.POST("/upload", middleware.RoleMiddleware("author", "admin"), deps.MediaCtrl.UploadFile)

			// Все аутентифицированные пользователи могут просматривать файлы
			protected.GET("", middleware.RoleMiddleware("user", "author", "moderator", "admin"), deps.MediaCtrl.GetAllMedia)
			protected.GET("/:id", middleware.RoleMiddleware("user", "author", "moderator", "admin"), deps.MediaCtrl.GetAllByArticleID)

			// Авторы могут удалять только свои файлы, администраторы — любые
			protected.DELETE("/:id", middleware.OwnershipMiddleware(), middleware.RoleMiddleware("author", "admin"), deps.MediaCtrl.DeleteFile)
		}
	}
}
