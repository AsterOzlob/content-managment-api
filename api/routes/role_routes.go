package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(router *gin.Engine, deps *Dependencies) {
	roleGroup := router.Group("/roles")
	{
		// Создание роли
		roleGroup.POST("", deps.RoleCtrl.CreateRole)

		// Получение всех ролей
		roleGroup.GET("", deps.RoleCtrl.GetAllRoles)

		// Получение роли по ID
		roleGroup.GET("/:id", deps.RoleCtrl.GetRoleByID)

		// Обновление роли
		roleGroup.PUT("/:id", deps.RoleCtrl.UpdateRole)

		// Удаление роли
		roleGroup.DELETE("/:id", deps.RoleCtrl.DeleteRole)
	}
}
