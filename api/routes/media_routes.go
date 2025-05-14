package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/middleware"
	"github.com/AsterOzlob/content_managment_api/pkg/appinit"
	"github.com/gin-gonic/gin"
)

// RegisterMediaRoutes настраивает маршруты для управления медиафайлами.
func RegisterMediaRoutes(r *gin.Engine, deps *appinit.Dependencies) {
	mediaGroup := r.Group("/media")
	{
		protected := mediaGroup.Group("/")
		protected.Use(middleware.AuthMiddleware(deps.JWTConfig)) // Middleware для JWT-аутентификации
		{
			// Только авторы могут загружать файлы
			protected.POST("/upload", middleware.RoleMiddleware("author", "admin"),
				deps.Controllers.MediaCtrl.UploadFileWithArticle)
			protected.POST("/upload/unlinked", middleware.RoleMiddleware("author", "admin"),
				deps.Controllers.MediaCtrl.UploadUnlinkedFile)

			// Все аутентифицированные пользователи могут просматривать
			protected.GET("", middleware.RoleMiddleware("user", "author", "moderator", "admin"),
				deps.Controllers.MediaCtrl.GetAllMedia)
			protected.GET("/:id", middleware.RoleMiddleware("user", "author", "moderator", "admin"),
				deps.Controllers.MediaCtrl.GetAllByArticleID)

			// Удаление только для владельцев или админов
			protected.DELETE("/:id", middleware.RoleMiddleware("author", "moderator", "admin"),
				deps.Controllers.MediaCtrl.DeleteFile)
		}
	}
}
