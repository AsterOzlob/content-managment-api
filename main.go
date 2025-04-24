package main

// @title           Content Management API
// @version         1.0
// @description     This is a RESTful API for managing users, roles, and content.
// @termsOfService  http://example.com/terms/

// @contact.name   Тимошенко Антон
// @contact.email  tumoshenko204@mail.ru

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer <your-token>

import (
	"github.com/AsterOzlob/content_managment_api/api/controllers"
	"github.com/AsterOzlob/content_managment_api/api/routes"
	"github.com/AsterOzlob/content_managment_api/config"
	_ "github.com/AsterOzlob/content_managment_api/docs"
	"github.com/AsterOzlob/content_managment_api/internal/repositories"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func main() {
	// Создаем экземпляр логгера для основного приложения.
	appLogger := logging.NewLogger("logs/app.log")
	appLogger.Log(logrus.InfoLevel, "Starting application", nil)

	// Инициализация приложения: загрузка конфигурации, подключение к базе данных и миграции.
	cfg, dbConn := initializeApp(appLogger)
	if cfg == nil || dbConn == nil {
		return
	}

	// Настройка зависимостей: репозитории, сервисы, контроллеры.
	deps := setupDependencies(dbConn, appLogger)

	// Настройка маршрутизатора и эндпоинтов API.
	r := setupRouter(deps, appLogger)

	// Определяем адрес сервера и запускаем HTTP-сервер.
	serverAddress := ":8080"
	appLogger.Log(logrus.InfoLevel, "Starting HTTP server", map[string]interface{}{
		"address": serverAddress,
	})
	if err := r.Run(serverAddress); err != nil {
		appLogger.Log(logrus.ErrorLevel, "Failed to start HTTP server", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	appLogger.Log(logrus.InfoLevel, "Application started successfully!", nil)
}

// initializeApp выполняет начальную настройку приложения:
// загрузку конфигурации, подключение к базе данных и выполнение миграций.
func initializeApp(logger *logging.Logger) (*config.Config, *gorm.DB) {
	// Загрузка конфигурации из .env файла или переменных окружения.
	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Log(logrus.ErrorLevel, "Error loading config", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, nil
	}

	// Инициализация подключения к базе данных PostgreSQL через GORM.
	dbConn, err := config.InitDB(cfg.DBConfig, logger)
	if err != nil {
		logger.Log(logrus.ErrorLevel, "Error initializing database connection", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, nil
	}

	// Выполнение миграций для создания таблиц в базе данных.
	if err := config.MigrateModels(dbConn, logger); err != nil {
		logger.Log(logrus.ErrorLevel, "Error migrating models", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, nil
	}

	return cfg, dbConn // Возвращаем конфигурацию и подключение к базе данных.
}

// setupDependencies настраивает зависимости приложения:
// репозитории, сервисы и контроллеры.
func setupDependencies(dbConn *gorm.DB, logger *logging.Logger) *routes.Dependencies {
	// Инициализация репозиториев.
	userRepo := repositories.NewUserRepository(dbConn, logger)
	articleRepo := repositories.NewArticleRepository(dbConn, logger)

	// Инициализация сервисов.
	userSvc := services.NewUserService(userRepo, logger)
	articleSvc := services.NewArticleService(articleRepo, logger)

	// Инциализация контроллеров.
	userCtrl := controllers.NewUserController(userSvc, logger)
	articleCtrl := controllers.NewArticleController(articleSvc, logger)

	// Возвращаем структуру зависимостей, которая содержит все компоненты приложения.
	return &routes.Dependencies{
		UserCtrl:    userCtrl,
		ArticleCtrl: articleCtrl,
		JWTConfig:   nil, // TODO: передать JWTConfig (например, из конфигурации).
	}
}

// setupRouter настраивает маршрутизатор Gin и определяет эндпоинты API.
func setupRouter(deps *routes.Dependencies, logger *logging.Logger) *gin.Engine {
	// Создаем новый экземпляр маршрутизатора Gin.
	r := gin.Default()

	// Регистрируем маршруты API, используя зависимости (контроллеры, сервисы и т.д.).
	routes.SetupRoutes(r, deps)

	// Добавляем эндпоинт для просмотра Swagger-документации.
	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.InstanceName("swagger"),
		ginSwagger.DocExpansion("none"), // Отключаем автоматическое раскрытие всех секций документации.
	))
	logger.Log(logrus.InfoLevel, "Swagger documentation endpoint configured", nil)

	return r // Возвращаем настроенный маршрутизатор.
}
