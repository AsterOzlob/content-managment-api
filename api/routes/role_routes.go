package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(router *gin.Engine, deps *Dependencies) {
	roleGroup := router.Group("/roles")
	{
		// Защищенные эндпоинты
		protected := roleGroup.Group("/")
		protected.Use(middleware.AuthMiddleware(deps.JWTConfig)) // Middleware для JWT-аутентификации
		protected.Use(middleware.RoleMiddleware("admin"))        // Только администраторы имеют доступ
		{
			// Создание роли
			protected.POST("", deps.RoleCtrl.CreateRole)

			// Получение всех ролей
			protected.GET("", deps.RoleCtrl.GetAllRoles)

			// Получение роли по ID
			protected.GET("/:id", deps.RoleCtrl.GetRoleByID)

			// Обновление роли
			protected.PUT("/:id", deps.RoleCtrl.UpdateRole)

			// Удаление роли
			protected.DELETE("/:id", deps.RoleCtrl.DeleteRole)
		}
	}
}
