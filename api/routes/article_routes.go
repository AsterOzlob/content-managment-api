package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/middleware"
	"github.com/AsterOzlob/content_managment_api/pkg/appinit"
	"github.com/gin-gonic/gin"
)

// RegisterArticleRoutes регистрирует маршруты для управления контентом.
func RegisterArticleRoutes(r *gin.Engine, deps *appinit.Dependencies) {
	content := r.Group("/articles")
	{
		// Открытые эндпоинты (без аутентификации)
		content.GET("", deps.Controllers.ArticleCtrl.GetAllArticles)     // Получение списка всех статей
		content.GET("/:id", deps.Controllers.ArticleCtrl.GetArticleByID) // Получение конкретной статьи

		// Защищенные эндпоинты
		protected := content.Group("/")
		protected.Use(middleware.AuthMiddleware(deps.JWTConfig)) // Middleware для JWT-аутентификации
		{
			// Авторы могут создавать статьи
			protected.POST("", middleware.RoleMiddleware("author", "admin"),
				deps.Controllers.ArticleCtrl.CreateArticle)

			// Авторы могут редактировать свои статьи, модераторы и администраторы — любые
			protected.PUT("/:id", middleware.RoleMiddleware("author", "moderator", "admin"),
				deps.Controllers.ArticleCtrl.UpdateArticle)

			// Авторы могут удалять свои статьи, модераторы и администраторы — любые
			protected.DELETE("/:id", middleware.RoleMiddleware("author", "moderator", "admin"),
				deps.Controllers.ArticleCtrl.DeleteArticle)
		}
	}
}
