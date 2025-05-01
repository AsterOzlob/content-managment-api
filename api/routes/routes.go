package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/controllers"
	"github.com/AsterOzlob/content_managment_api/config"
	"github.com/AsterOzlob/content_managment_api/internal/database/repositories"

	"github.com/gin-gonic/gin"
)

// Dependencies содержит зависимости для маршрутов
type Dependencies struct {
	AuthCtrl         *controllers.AuthController
	UserCtrl         *controllers.UserController
	ArticleCtrl      *controllers.ArticleController
	CommentCtrl      *controllers.CommentController
	MediaCtrl        *controllers.MediaController
	RoleCtrl         *controllers.RoleController
	JWTConfig        *config.JWTConfig
	RefreshTokenRepo *repositories.RefreshTokenRepository
}

func SetupRoutes(router *gin.Engine, deps *Dependencies) {
	// Регистрация маршрутов для аутентификации
	RegisterAuthRoutes(router, deps)
	// Регистрация маршрутов для пользователей
	RegisterUserRoutes(router, deps)
	// Регистрация маршрутов для контента
	RegisterArticleRoutes(router, deps)
	// Регистрация маршрутов для комментариев
	RegisterCommentRoutes(router, deps)
	// Регистрация маршрутов для медиа файлов
	RegisterMediaRoutes(router, deps)
	// Регистрация маршрутов для медиа ролей
	RegisterRoleRoutes(router, deps)
}
