package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/middleware"
	"github.com/AsterOzlob/content_managment_api/cmd/app"
	"github.com/gin-gonic/gin"
)

// RegisterRoleRoutes регистрирует маршруты для управления ролями.
func RegisterRoleRoutes(router *gin.Engine, deps *app.Dependencies) {
	roleGroup := router.Group("/roles")
	{
		// Защищенные эндпоинты
		protected := roleGroup.Group("/")
		protected.Use(middleware.AuthMiddleware(deps.JWTConfig)) // Middleware для JWT-аутентификации
		protected.Use(middleware.RoleMiddleware("admin"))        // Только администраторы имеют доступ
		{
			// Создание роли
			protected.POST("", deps.Controllers.RoleCtrl.CreateRole)

			// Получение всех ролей
			protected.GET("", deps.Controllers.RoleCtrl.GetAllRoles)

			// Получение роли по ID
			protected.GET("/:id", deps.Controllers.RoleCtrl.GetRoleByID)

			// Обновление роли
			protected.PUT("/:id", deps.Controllers.RoleCtrl.UpdateRole)

			// Удаление роли
			protected.DELETE("/:id", deps.Controllers.RoleCtrl.DeleteRole)
		}
	}
}
