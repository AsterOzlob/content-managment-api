package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/controllers"
	"github.com/AsterOzlob/content_managment_api/config"
	"github.com/gin-gonic/gin"
)

// Dependencies содержит зависимости для маршрутов
type Dependencies struct {
	UserCtrl    *controllers.UserController
	ArticleCtrl *controllers.ArticleController
	CommentCtrl *controllers.CommentController
	JWTConfig   *config.JWTConfig
}

func SetupRoutes(router *gin.Engine, deps *Dependencies) {
	// Регистрация маршрутов для пользователей
	RegisterUserRoutes(router, deps)
	// Регистрация маршрутов для контента
	RegisterArticleRoutes(router, deps)
	// Регистрация маршрутов для комментариев
	RegisterCommentRoutes(router, deps)
}
