package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/middleware"
	"github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/AsterOzlob/content_managment_api/pkg/appinit"
	"github.com/gin-gonic/gin"
)

// setupRouter настраивает маршрутизатор Gin и определяет эндпоинты API.
func SetupRouter(deps *appinit.Dependencies, logger logger.Logger) *gin.Engine {
	r := gin.Default()

	// Применяем middleware
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.XSSMiddleware())
	r.Use(middleware.NewRateLimiter(100).Middleware())

	// Регистрируем маршруты
	SetupRoutes(r, deps)

	return r
}
