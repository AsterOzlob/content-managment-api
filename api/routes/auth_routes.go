package routes

import (
	"github.com/AsterOzlob/content_managment_api/cmd/app"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes регистрирует маршруты для аутентификации.
func RegisterAuthRoutes(r *gin.Engine, deps *app.Dependencies) {
	auth := r.Group("/users")
	{
		// Открытые эндпоинты
		auth.POST("/signup", deps.Controllers.AuthCtrl.SignUp)        // Регистрация пользователя
		auth.POST("/login", deps.Controllers.AuthCtrl.Login)          // Вход пользователя
		auth.POST("/refresh", deps.Controllers.AuthCtrl.RefreshToken) // Обновление токенов
	}
}
