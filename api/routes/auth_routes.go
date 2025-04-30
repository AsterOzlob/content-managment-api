package routes

import (
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes регистрирует маршруты для аутентификации.
func RegisterAuthRoutes(r *gin.Engine, deps *Dependencies) {
	auth := r.Group("/auth")
	{
		// Открытые эндпоинты
		auth.POST("/signup", deps.AuthCtrl.SignUp)        // Регистрация пользователя
		auth.POST("/login", deps.AuthCtrl.Login)          // Вход пользователя
		auth.POST("/refresh", deps.AuthCtrl.RefreshToken) // Обновление токенов
	}
}
