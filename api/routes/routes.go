package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/controllers"
	"github.com/AsterOzlob/content_managment_api/api/middleware/auth"
	"github.com/gin-gonic/gin"
)

// Dependencies содержит зависимости для маршрутов
type Dependencies struct {
	UserCtrl   *controllers.UserController
	JWTManager *auth.JWTManager
}

func SetupRoutes(router *gin.Engine, deps *Dependencies) {
	// Регистрация маршрутов для пользователей
	RegisterUserRoutes(router, deps)
}
