package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterCommentRoutes регистрирует маршруты для управления комментариями.
func RegisterCommentRoutes(r *gin.Engine, deps *Dependencies) {
	// Группа маршрутов для комментариев, привязанных к статьям
	content := r.Group("/articles")
	{
		protected := content.Group("/")
		protected.Use(middleware.AuthMiddleware(deps.JWTConfig)) // Middleware для JWT-аутентификации
		{
			// Все аутентифицированные пользователи могут оставлять комментарии
			protected.POST("/:id/comments", middleware.RoleMiddleware("user", "author", "moderator", "admin"), deps.CommentCtrl.AddCommentToArticle)

			// Все аутентифицированные пользователи могут просматривать комментарии
			protected.GET("/:id/comments", middleware.RoleMiddleware("user", "author", "moderator", "admin"), deps.CommentCtrl.GetCommentsByArticleID)
		}
	}

	// Глобальные маршруты для комментариев
	r.DELETE("/comments/:id",
		middleware.AuthMiddleware(deps.JWTConfig),
		middleware.OwnershipMiddleware(),
		middleware.RoleMiddleware("user", "moderator", "admin"),
		deps.CommentCtrl.DeleteComment,
	)
}
