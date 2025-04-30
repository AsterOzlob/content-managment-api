package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes регистрирует маршруты для управления пользователями.
func RegisterUserRoutes(r *gin.Engine, deps *Dependencies) {
	user := r.Group("/users")
	{
		// Защищенные эндпоинты
		protected := user.Group("/")
		protected.Use(middleware.AuthMiddleware(deps.JWTConfig)) // Middleware для JWT-аутентификации
		{
			// Временно убрана проверка прав (middleware.RoleMiddleware)
			protected.GET("/:id", deps.UserCtrl.GetUserByID)
			protected.PATCH("/:id/role", deps.UserCtrl.AssignRole)
			protected.DELETE("/:id", deps.UserCtrl.DeleteUser)
			protected.GET("", deps.UserCtrl.GetAllUsers)
		}
	}
}
