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
			// Пользователь может получить информацию только о себе
			protected.GET("/:id", middleware.OwnershipMiddleware(), middleware.RoleMiddleware("user", "admin"), deps.UserCtrl.GetUserByID)

			// Администраторы могут назначать роли
			protected.PATCH("/:id/role", middleware.RoleMiddleware("admin"), deps.UserCtrl.AssignRole)

			// Администраторы могут удалять пользователей
			protected.DELETE("/:id", middleware.RoleMiddleware("admin"), deps.UserCtrl.DeleteUser)

			// Администраторы могут получать список всех пользователей
			protected.GET("", middleware.RoleMiddleware("admin"), deps.UserCtrl.GetAllUsers)
		}
	}
}
