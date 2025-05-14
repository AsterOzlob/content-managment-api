package routes

import (
	"github.com/AsterOzlob/content_managment_api/cmd/app"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, deps *app.Dependencies) {
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
