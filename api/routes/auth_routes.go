package routes

import (
	"github.com/AsterOzlob/content_managment_api/pkg/appinit"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes регистрирует маршруты для аутентификации.
func RegisterAuthRoutes(r *gin.Engine, deps *appinit.Dependencies) {
	auth := r.Group("/users")
	{
		// Открытые эндпоинты
		auth.POST("/signup", deps.Controllers.AuthCtrl.SignUp) // Регистрация пользователя
		auth.POST("/login", deps.Controllers.AuthCtrl.Login)   // Вход пользователя
	}
}
