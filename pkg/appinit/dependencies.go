package appinit

import (
	"github.com/AsterOzlob/content_managment_api/api/controllers"
	"github.com/AsterOzlob/content_managment_api/config"
	"github.com/AsterOzlob/content_managment_api/internal/database/repositories"
	"github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	"gorm.io/gorm"
)

// Loggers содержит все логгеры проекта
type Loggers struct {
	AuthLogger    logger.Logger
	UserLogger    logger.Logger
	ArticleLogger logger.Logger
	CommentLogger logger.Logger
	MediaLogger   logger.Logger
	RoleLogger    logger.Logger
}

// Repositories содержит все репозитории проекта
type Repositories struct {
	UserRepo         *repositories.UserRepository
	ArticleRepo      *repositories.ArticleRepository
	CommentRepo      *repositories.CommentRepository
	MediaRepo        *repositories.MediaRepository
	RefreshTokenRepo *repositories.RefreshTokenRepository
	RoleRepo         *repositories.RoleRepository
}

// Services содержит все сервисы проекта
type Services struct {
	AuthService    *services.AuthService
	UserService    *services.UserService
	ArticleService *services.ArticleService
	CommentService *services.CommentService
	MediaService   *services.MediaService
	RoleService    *services.RoleService
}

// Controllers содержит все контроллеры проекта
type Controllers struct {
	AuthCtrl    *controllers.AuthController
	UserCtrl    *controllers.UserController
	ArticleCtrl *controllers.ArticleController
	CommentCtrl *controllers.CommentController
	MediaCtrl   *controllers.MediaController
	RoleCtrl    *controllers.RoleController
}

// Dependencies содержит все зависимости проекта
type Dependencies struct {
	Controllers  *Controllers
	Services     *Services
	Repositories *Repositories
	Loggers      *Loggers
	JWTConfig    *config.JWTConfig
	MediaConfig  *config.MediaConfig
}

// SetupDependencies настраивает зависимости приложения:
// логгеры, репозитории, сервисы и контроллеры.
func SetupDependencies(dbConn *gorm.DB, cfg *config.Config) *Dependencies {
	// Инициализация логгеров
	loggers := setupLoggers()

	// Инициализация репозиториев
	repos := setupRepositories(dbConn, loggers)

	// Инициализация сервисов
	services := setupServices(repos, cfg, loggers)

	// Инициализация контроллеров
	controllers := setupControllers(services, cfg)

	// Возвращаем структуру зависимостей
	return &Dependencies{
		Controllers:  controllers,
		Services:     services,
		Repositories: repos,
		Loggers:      loggers,
		JWTConfig:    cfg.JWTConfig,
		MediaConfig:  cfg.MediaConfig,
	}
}

// setupLoggers создает логгеры для каждой области
func setupLoggers() *Loggers {
	return &Loggers{
		AuthLogger:    logger.NewLogger("logs/auth.log"),
		UserLogger:    logger.NewLogger("logs/users.log"),
		ArticleLogger: logger.NewLogger("logs/articles.log"),
		CommentLogger: logger.NewLogger("logs/comments.log"),
		MediaLogger:   logger.NewLogger("logs/media.log"),
		RoleLogger:    logger.NewLogger("logs/roles.log"),
	}
}

// setupRepositories инициализирует репозитории
func setupRepositories(dbConn *gorm.DB, loggers *Loggers) *Repositories {
	return &Repositories{
		UserRepo:         repositories.NewUserRepository(dbConn, loggers.UserLogger),
		ArticleRepo:      repositories.NewArticleRepository(dbConn, loggers.ArticleLogger),
		CommentRepo:      repositories.NewCommentRepository(dbConn, loggers.CommentLogger),
		MediaRepo:        repositories.NewMediaRepository(dbConn, loggers.MediaLogger),
		RefreshTokenRepo: repositories.NewRefreshTokenRepository(dbConn, loggers.AuthLogger),
		RoleRepo:         repositories.NewRoleRepository(dbConn, loggers.RoleLogger),
	}
}

// setupServices инициализирует сервисы
func setupServices(repos *Repositories, cfg *config.Config, loggers *Loggers) *Services {
	return &Services{
		AuthService: services.NewAuthService(
			repos.UserRepo,
			repos.RefreshTokenRepo,
			loggers.AuthLogger,
			cfg.JWTConfig,
		),
		UserService: services.NewUserService(
			repos.UserRepo,
			loggers.UserLogger,
		),
		ArticleService: services.NewArticleService(
			repos.ArticleRepo,
			loggers.ArticleLogger,
		),
		CommentService: services.NewCommentService(
			repos.CommentRepo,
			loggers.CommentLogger,
		),
		MediaService: services.NewMediaService(
			repos.MediaRepo,
			repos.ArticleRepo,
			loggers.MediaLogger,
		),
		RoleService: services.NewRoleService(
			repos.RoleRepo,
			loggers.RoleLogger,
		),
	}
}

// setupControllers инициализирует контроллеры
func setupControllers(services *Services, cfg *config.Config) *Controllers {
	return &Controllers{
		AuthCtrl:    controllers.NewAuthController(services.AuthService),
		UserCtrl:    controllers.NewUserController(services.UserService),
		ArticleCtrl: controllers.NewArticleController(services.ArticleService),
		CommentCtrl: controllers.NewCommentController(services.CommentService),
		MediaCtrl: controllers.NewMediaController(
			services.MediaService,
			cfg.MediaConfig,
		),
		RoleCtrl: controllers.NewRoleController(services.RoleService),
	}
}
