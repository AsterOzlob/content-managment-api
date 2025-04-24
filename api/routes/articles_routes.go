package routes

import (
	"github.com/gin-gonic/gin"
)

// RegisterContentRoutes регистрирует маршруты для управления контентом.
func RegisterContentRoutes(r *gin.Engine, deps *Dependencies) {
	content := r.Group("/article")
	{
		// Открытые эндпоинты (без аутентификации)
		content.GET("", deps.ArticleCtrl.GetAllArticles)     // Получение списка всех статей и новостей
		content.GET("/:id", deps.ArticleCtrl.GetArticleByID) // Получение конкретного контента по ID

		// Защищенные эндпоинты (временно без аутентификации)
		protected := content.Group("/")
		{
			protected.POST("", deps.ArticleCtrl.CreateArticle)       // Создание нового контента
			protected.PUT("/:id", deps.ArticleCtrl.UpdateArticle)    // Обновление существующего контента
			protected.DELETE("/:id", deps.ArticleCtrl.DeleteArticle) // Удаление контента
			// protected.POST("/:id/comments", deps.ContentCtrl.AddCommentToContent) // Добавление комментария к контенту
		}
	}
}
