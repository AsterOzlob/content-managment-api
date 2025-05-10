package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterCommentRoutes регистрирует маршруты для управления комментариями.
func RegisterCommentRoutes(r *gin.Engine, deps *Dependencies) {
	content := r.Group("/articles")
	{
		protected := content.Group("/")
		protected.Use(middleware.AuthMiddleware(deps.JWTConfig))
		{
			protected.POST("/:id/comments", middleware.RoleMiddleware("user", "author", "moderator", "admin"), deps.CommentCtrl.AddCommentToArticle)
			protected.GET("/:id/comments", middleware.RoleMiddleware("user", "author", "moderator", "admin"), deps.CommentCtrl.GetCommentsByArticleID)

			// Добавляем маршрут для редактирования комментария
			protected.PUT("/comments/:id", middleware.RoleMiddleware("user", "author", "moderator", "admin"), deps.CommentCtrl.UpdateComment)
		}
	}

	r.DELETE("/comments/:id",
		middleware.AuthMiddleware(deps.JWTConfig),
		middleware.RoleMiddleware("user", "author", "moderator", "admin"),
		deps.CommentCtrl.DeleteComment,
	)
}
