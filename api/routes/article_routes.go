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
			// Авторы могут создавать и управлять своими материалами
			protected.POST("", middleware.RoleMiddleware("author", "admin"), deps.ArticleCtrl.CreateArticle)       // Создание статьи
			protected.PUT("/:id", middleware.RoleMiddleware("author", "admin"), deps.ArticleCtrl.UpdateArticle)    // Обновление статьи
			protected.DELETE("/:id", middleware.RoleMiddleware("author", "admin"), deps.ArticleCtrl.DeleteArticle) // Удаление статьи
		}
	}
}
