package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/AsterOzlob/content_managment_api/internal/models"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/sirupsen/logrus"
)

// InitDB инициализирует подключение к базе данных.
func InitDB(dbConfig *DBConfig, logger *logging.Logger) (*gorm.DB, error) {
	// Формирование строки подключения (DSN)
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		dbConfig.DBHost,
		dbConfig.DBPort,
		dbConfig.DBUser,
		dbConfig.DBName,
		getEnv("DB_SSL_MODE", "disable"), // Значение по умолчанию: disable
		dbConfig.DBPassword,
	)

	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log(logrus.ErrorLevel, "Failed to connect to database", map[string]interface{}{
			"dsn":   dsn,
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверка соединения
	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		logger.Log(logrus.ErrorLevel, "Failed to ping database", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Log(logrus.InfoLevel, "Database connection established successfully!", nil)
	return db, nil
}

// MigrateModels выполняет миграцию моделей.
func MigrateModels(db *gorm.DB, logger *logging.Logger) error {
	models := []interface{}{
		&models.User{},
		&models.Role{},
		&models.RefreshToken{},
		&models.Article{},
		&models.Media{},
		&models.Comment{},
	}

	// Миграция моделей
	if err := db.AutoMigrate(models...); err != nil {
		logger.Log(logrus.ErrorLevel, "Failed to migrate models", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to migrate models: %w", err)
	}

	logger.Log(logrus.InfoLevel, "Database migrations completed successfully!", nil)
	return nil
}
