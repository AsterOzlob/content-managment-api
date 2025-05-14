package routes

import (
	"github.com/AsterOzlob/content_managment_api/api/middleware"
	"github.com/AsterOzlob/content_managment_api/cmd/app"
	_ "github.com/AsterOzlob/content_managment_api/docs"
	"github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// setupRouter настраивает маршрутизатор Gin и определяет эндпоинты API.
func SetupRouter(deps *app.Dependencies, logger logger.Logger) *gin.Engine {
	r := gin.Default()

	// Применяем middleware
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.XSSMiddleware())
	r.Use(middleware.NewRateLimiter(100).Middleware())

	// Регистрируем маршруты
	SetupRoutes(r, deps)

	// Добавляем эндпоинт для Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
