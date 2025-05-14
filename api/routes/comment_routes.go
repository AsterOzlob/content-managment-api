package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/middleware"
	"github.com/AsterOzlob/content_managment_api/pkg/appinit"
	"github.com/gin-gonic/gin"
)

// RegisterCommentRoutes регистрирует маршруты для управления комментариями.
func RegisterCommentRoutes(r *gin.Engine, deps *appinit.Dependencies) {
	content := r.Group("/articles")
	{
		protected := content.Group("/")
		protected.Use(middleware.AuthMiddleware(deps.JWTConfig)) // Middleware для JWT-аутентификации
		{
			// Добавление комментария к статье
			protected.POST("/:id/comments", middleware.RoleMiddleware("user", "author", "moderator", "admin"),
				deps.Controllers.CommentCtrl.AddCommentToArticle)

			// Получение комментариев для статьи
			protected.GET("/:id/comments", middleware.RoleMiddleware("user", "author", "moderator", "admin"),
				deps.Controllers.CommentCtrl.GetCommentsByArticleID)

			// Редактирование комментария
			protected.PUT("/comments/:id", middleware.RoleMiddleware("user", "author", "moderator", "admin"),
				deps.Controllers.CommentCtrl.UpdateComment)
		}
	}

	// Удаление комментария
	r.DELETE("/comments/:id",
		middleware.AuthMiddleware(deps.JWTConfig),
		middleware.RoleMiddleware("user", "author", "moderator", "admin"),
		deps.Controllers.CommentCtrl.DeleteComment,
	)
}
