package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, deps *Dependencies) {
	// Открытые эндпоинты
	r.POST("/users/signup", deps.UserCtrl.SignUp)
	r.POST("/users/login", deps.UserCtrl.Login)

	// Защищённые эндпоинты
	r.GET("/users/:id", deps.UserCtrl.GetUserByID)
	r.PATCH("/users/:id/role", deps.UserCtrl.AssignRole)
	r.DELETE("/users/:id", deps.UserCtrl.DeleteUser)
	r.GET("/users", deps.UserCtrl.GetAllUsers)
}
