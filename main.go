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
	"time"

	"github.com/AsterOzlob/content_managment_api/api/controllers"
	"github.com/AsterOzlob/content_managment_api/api/routes"
	"github.com/AsterOzlob/content_managment_api/config"
	_ "github.com/AsterOzlob/content_managment_api/docs"
	"github.com/AsterOzlob/content_managment_api/internal/database"
	"github.com/AsterOzlob/content_managment_api/internal/database/repositories"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func main() {
	// Создаем экземпляр логгера для основного приложения.
	appLogger := logger.NewLogger("logs/app.log")
	appLogger.Info("Starting application")

	// Инициализация приложения: загрузка конфигурации, подключение к базе данных и миграции.
	cfg, dbConn := initializeApp(appLogger)
	if cfg == nil || dbConn == nil {
		appLogger.Error("Failed to initialize application")
		return
	}

	// Настройка зависимостей: репозитории, сервисы, контроллеры.
	deps := setupDependencies(dbConn, cfg)

	// Запуск планировщика для очистки истекших токенов.
	startTokenCleanupScheduler(deps.RefreshTokenRepo, appLogger)

	// Настройка маршрутизатора и эндпоинтов API.
	r := setupRouter(deps, appLogger)

	// Определяем адрес сервера и запускаем HTTP-сервер.
	serverAddress := ":8080"
	appLogger.WithField("address", serverAddress).Info("Starting HTTP server")
	if err := r.Run(serverAddress); err != nil {
		appLogger.WithError(err).Error("Failed to start HTTP server")
		return
	}

	appLogger.Info("Application started successfully!")
}

// initializeApp выполняет начальную настройку приложения:
// загрузку конфигурации, подключение к базе данных и выполнение миграций.
func initializeApp(logger logger.Logger) (*config.Config, *gorm.DB) {
	// Загрузка конфигурации из .env файла или переменных окружения.
	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.WithError(err).Error("Error loading config")
		return nil, nil
	}

	// Инициализация подключения к базе данных PostgreSQL через GORM.
	dbConn, err := database.InitDB(cfg.DBConfig, logger)
	if err != nil {
		logger.WithError(err).Error("Error initializing database connection")
		return nil, nil
	}

	// Выполнение миграций для создания таблиц в базе данных.
	if err := database.MigrateModels(dbConn, logger); err != nil {
		logger.WithError(err).Error("Error migrating models")
		return nil, nil
	}

	return cfg, dbConn // Возвращаем конфигурацию и подключение к базе данных.
}

// setupDependencies настраивает зависимости приложения:
// репозитории, сервисы и контроллеры.
func setupDependencies(dbConn *gorm.DB, cfg *config.Config) *routes.Dependencies {
	// Создаем отдельные логгеры для каждой области.
	authLogger := logger.NewLogger("logs/auth.log")
	userLogger := logger.NewLogger("logs/users.log")
	articleLogger := logger.NewLogger("logs/articles.log")
	commentLogger := logger.NewLogger("logs/comments.log")
	mediaLogger := logger.NewLogger("logs/media.log")
	roleLogger := logger.NewLogger("logs/roles.log")

	// Инициализация репозиториев.
	userRepo := repositories.NewUserRepository(dbConn, userLogger)
	articleRepo := repositories.NewArticleRepository(dbConn, articleLogger)
	commentRepo := repositories.NewCommentRepository(dbConn, commentLogger)
	mediaRepo := repositories.NewMediaRepository(dbConn, mediaLogger)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(dbConn, authLogger)
	roleRepo := repositories.NewRoleRepository(dbConn, roleLogger)

	// Инициализация сервисов.
	authService := services.NewAuthService(userRepo, refreshTokenRepo, authLogger, cfg.JWTConfig)
	userSvc := services.NewUserService(userRepo, userLogger)
	articleSvc := services.NewArticleService(articleRepo, articleLogger)
	commentSvc := services.NewCommentService(commentRepo, commentLogger)
	mediaSvc := services.NewMediaService(mediaRepo, mediaLogger)
	roleSvc := services.NewRoleService(roleRepo, roleLogger)

	// Инициализация контроллеров.
	authCtrl := controllers.NewAuthController(authService, authLogger)
	userCtrl := controllers.NewUserController(userSvc, userLogger)
	articleCtrl := controllers.NewArticleController(articleSvc, articleLogger)
	commentCtrl := controllers.NewCommentController(commentSvc, commentLogger)
	mediaCtrl := controllers.NewMediaController(mediaSvc, mediaLogger, cfg.MediaConfig)
	roleCtrl := controllers.NewRoleController(roleSvc, roleLogger)

	// Возвращаем структуру зависимостей, которая содержит все компоненты приложения.
	return &routes.Dependencies{
		AuthCtrl:         authCtrl,
		UserCtrl:         userCtrl,
		ArticleCtrl:      articleCtrl,
		CommentCtrl:      commentCtrl,
		MediaCtrl:        mediaCtrl,
		RoleCtrl:         roleCtrl,
		JWTConfig:        cfg.JWTConfig,
		RefreshTokenRepo: refreshTokenRepo,
	}
}

// setupRouter настраивает маршрутизатор Gin и определяет эндпоинты API.
func setupRouter(deps *routes.Dependencies, logger logger.Logger) *gin.Engine {
	// Создаем новый экземпляр маршрутизатора Gin.
	r := gin.Default()

	// Регистрируем маршруты API, используя зависимости (контроллеры, сервисы и т.д.).
	routes.SetupRoutes(r, deps)

	// Добавляем эндпоинт для просмотра Swagger-документации.
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logger.Info("Swagger documentation endpoint configured")

	return r
}

// startTokenCleanupScheduler запускает планировщик для очистки истекших токенов.
func startTokenCleanupScheduler(refreshTokenRepo *repositories.RefreshTokenRepository, logger logger.Logger) {
	go func() {
		for {
			time.Sleep(1 * time.Hour) // Запуск каждый час
			logger.Info("Running scheduled cleanup of expired refresh tokens")
			if err := refreshTokenRepo.CleanupExpiredTokens(); err != nil {
				logger.WithError(err).Error("Error during cleanup of expired refresh tokens")
			} else {
				logger.Info("Successfully cleaned up expired refresh tokens")
			}
		}
	}()
}
