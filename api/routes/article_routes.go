package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterArticleRoutes регистрирует маршруты для управления контентом.
func RegisterArticleRoutes(r *gin.Engine, deps *Dependencies) {
	content := r.Group("/articles")
	{
		// Открытые эндпоинты (без аутентификации)
		content.GET("", deps.ArticleCtrl.GetAllArticles)     // Получение списка всех статей
		content.GET("/:id", deps.ArticleCtrl.GetArticleByID) // Получение конкретной статьи

		// Защищенные эндпоинты
		protected := content.Group("/")
		protected.Use(middleware.AuthMiddleware(deps.JWTConfig)) // Middleware для JWT-аутентификации
		{
			// Авторы могут создавать статьи
			protected.POST("", middleware.RoleMiddleware("author", "admin"), deps.ArticleCtrl.CreateArticle)

			// Авторы могут редактировать свои статьи, модераторы и администраторы — любые
			protected.PUT("/:id", middleware.OwnershipMiddleware(), middleware.RoleMiddleware("author", "moderator", "admin"), deps.ArticleCtrl.UpdateArticle)

			// Авторы могут удалять свои статьи, модераторы и администраторы — любые
			protected.DELETE("/:id", middleware.OwnershipMiddleware(), middleware.RoleMiddleware("author", "moderator", "admin"), deps.ArticleCtrl.DeleteArticle)
		}
	}
}
