package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/middleware"
	"github.com/AsterOzlob/content_managment_api/pkg/appinit"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes регистрирует маршруты для управления пользователями.
func RegisterUserRoutes(r *gin.Engine, deps *appinit.Dependencies) {
	user := r.Group("/users")
	{
		// Защищенные эндпоинты
		protected := user.Group("/")
		protected.Use(middleware.AuthMiddleware(deps.JWTConfig)) // Middleware для JWT-аутентификации
		{
			// Пользователь может получить информацию только о себе
			protected.GET("/:id", middleware.RoleMiddleware("user", "admin"),
				deps.Controllers.UserCtrl.GetUserByID)

			// Администраторы могут назначать роли
			protected.PATCH("/:id/role", middleware.RoleMiddleware("admin"),
				deps.Controllers.UserCtrl.AssignRole)

			// Администраторы могут удалять пользователей
			protected.DELETE("/:id", middleware.RoleMiddleware("admin"),
				deps.Controllers.UserCtrl.DeleteUser)

			// Администраторы могут получать список всех пользователей
			protected.GET("", middleware.RoleMiddleware("admin"),
				deps.Controllers.UserCtrl.GetAllUsers)
		}
	}
}
