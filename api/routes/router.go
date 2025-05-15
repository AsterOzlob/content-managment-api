package routes

import (
	"net/http"

	"github.com/AsterOzlob/content_managment_api/api/middleware"
	"github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/AsterOzlob/content_managment_api/pkg/appinit"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter настраивает маршрутизатор Gin и определяет эндпоинты API.
func SetupRouter(deps *appinit.Dependencies, logger logger.Logger) *gin.Engine {
	r := gin.Default()

	// Подключаем Swagger UI с кастомной настройкой
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.InstanceName("swagger"),
		ginSwagger.DocExpansion("none"),
	))

	// Редирект /docs → /swagger/index.html
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Применяем middleware
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.NewRateLimiter(100).Middleware())

	// Регистрируем маршруты API
	SetupRoutes(r, deps)

	return r
}
