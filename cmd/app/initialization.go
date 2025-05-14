package app

import (
	"github.com/AsterOzlob/content_managment_api/config"
	"github.com/AsterOzlob/content_managment_api/internal/database"
	"github.com/AsterOzlob/content_managment_api/internal/logger"
	"gorm.io/gorm"
)

// InitializeApp выполняет начальную настройку приложения:
// загрузку конфигурации, подключение к базе данных и выполнение миграций.
func InitializeApp(logger logger.Logger) (*config.Config, *gorm.DB) {
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
