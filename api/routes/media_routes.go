package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterMediaRoutes настраивает маршруты для управления медиафайлами.
func RegisterMediaRoutes(r *gin.Engine, deps *Dependencies) {
	mediaGroup := r.Group("/media")
	{
		protected := mediaGroup.Group("/")
		protected.Use(middleware.AuthMiddleware(deps.JWTConfig))
		{
			// Только авторы могут загружать файлы
			protected.POST("/upload", middleware.RoleMiddleware("author", "admin"), deps.MediaCtrl.UploadFileWithArticle)
			protected.POST("/upload/unlinked", middleware.RoleMiddleware("author", "admin"), deps.MediaCtrl.UploadUnlinkedFile)

			// Все аутентифицированные пользователи могут просматривать
			protected.GET("", middleware.RoleMiddleware("user", "author", "moderator", "admin"), deps.MediaCtrl.GetAllMedia)
			protected.GET("/:id", middleware.RoleMiddleware("user", "author", "moderator", "admin"), deps.MediaCtrl.GetAllByArticleID)

			// Удаление только для владельцев или админов
			protected.DELETE("/:id", middleware.RoleMiddleware("author", "moderator", "admin"), deps.MediaCtrl.DeleteFile)
		}
	}
}
