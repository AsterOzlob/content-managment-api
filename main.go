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
	"fmt"

	"github.com/AsterOzlob/content_managment_api/api/controllers"
	"github.com/AsterOzlob/content_managment_api/api/routes"
	"github.com/AsterOzlob/content_managment_api/config"
	_ "github.com/AsterOzlob/content_managment_api/docs"
	"github.com/AsterOzlob/content_managment_api/internal/repositories"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// Инициализация подключения к БД
	dbConn, err := config.InitDB(cfg.DBConfig)
	if err != nil {
		fmt.Println("Error initializing database connection:", err)
		return
	}

	// Миграция моделей
	err = config.MigrateModels(dbConn)
	if err != nil {
		fmt.Println("Error migrating models:", err)
		return
	}

	// DI для логгеров
	articleLogger := logging.NewLogger("logs/articles.log")
	userLogger := logging.NewLogger("logs/users.log")

	// DI для репозиториев
	userRepo := repositories.NewUserRepository(dbConn, userLogger)
	articleRepo := repositories.NewArticleRepository(dbConn, articleLogger)

	// DI для сервисов
	userSvc := services.NewUserService(userRepo, userLogger)
	articleSvc := services.NewArticleService(articleRepo, articleLogger)

	// DI для контроллеров
	userCtrl := controllers.NewUserController(userSvc, userLogger)
	articleCtrl := controllers.NewArticleController(articleSvc, articleLogger)

	// Инициализация структуры зависимостей
	deps := &routes.Dependencies{
		UserCtrl:    userCtrl,
		ArticleCtrl: articleCtrl,
		JWTConfig:   cfg.JWTConfig,
	}

	// Инициализация маршрутизатора
	r := gin.Default()

	// Настройка маршрутов
	routes.SetupRoutes(r, deps)

	// Найстройка Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.InstanceName("swagger"),
		ginSwagger.DocExpansion("none"),
	))

	// Запуск сервера
	r.Run(":8080")

	fmt.Println("Application started successfully!")
}
