package routes

import (
	"github.com/gin-gonic/gin"
)

// RegisterCommentRoutes регистрирует маршруты для управления комментариями.
func RegisterCommentRoutes(r *gin.Engine, deps *Dependencies) {
	// Группа маршрутов для комментариев, привязанных к статьям
	content := r.Group("/articles")
	{
		protected := content.Group("/")
		{
			protected.POST("/:id/comments", deps.CommentCtrl.AddCommentToArticle)   // Добавление комментария к статье
			protected.GET("/:id/comments", deps.CommentCtrl.GetCommentsByArticleID) // Получение всех комментариев к статье
		}
	}

	// Глобальные маршруты для комментариев
	r.DELETE("/comments/:id", deps.CommentCtrl.DeleteComment) // Удаление комментария
}
