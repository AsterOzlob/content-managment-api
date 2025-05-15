package main

import (
	"log"

	"github.com/AsterOzlob/content_managment_api/api/routes"
	_ "github.com/AsterOzlob/content_managment_api/docs"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/AsterOzlob/content_managment_api/internal/scheduler"
	"github.com/AsterOzlob/content_managment_api/pkg/appinit"
	"github.com/joho/godotenv"
)

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
func main() {
	// Создаем экземпляр логгера для основного приложения.
	appLogger := logger.NewLogger("logs/app.log")
	appLogger.Info("Starting application")

	// Загружаем .env файл, если он есть
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл:", err)
	}

	// Инициализация приложения: загрузка конфигурации, подключение к базе данных и миграции.
	cfg, dbConn := appinit.InitializeApp(appLogger)
	if cfg == nil || dbConn == nil {
		appLogger.Error("Failed to initialize application")
		return
	}

	// Настройка зависимостей: репозитории, сервисы, контроллеры.
	deps := appinit.SetupDependencies(dbConn, cfg)

	// Запуск планировщика для очистки истекших токенов.
	scheduler.StartTokenCleanupScheduler(deps.Repositories.RefreshTokenRepo, appLogger)

	// Настройка маршрутизатора и эндпоинтов API.
	r := routes.SetupRouter(deps, appLogger)

	// Определяем адрес сервера и запускаем HTTP-сервер.
	serverAddress := ":8080"
	if err := r.Run(serverAddress); err != nil {
		appLogger.WithError(err).Error("Failed to start HTTP server")
		return
	}

	appLogger.Info("Application started successfully!")
}
