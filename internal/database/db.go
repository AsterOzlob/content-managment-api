package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/AsterOzlob/content_managment_api/config"
	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/sirupsen/logrus"
)

// InitDB инициализирует подключение к базе данных.
func InitDB(dbConfig *config.DBConfig, logger logger.Logger) (*gorm.DB, error) {
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
		logger.WithFields(logrus.Fields{
			"dsn":   dsn,
			"error": err.Error(),
		}).Error("Failed to connect to database")
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверка соединения
	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		logger.WithError(err).Error("Failed to ping database")
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established successfully!")
	return db, nil
}

// MigrateModels выполняет миграцию моделей.
func MigrateModels(db *gorm.DB, logger logger.Logger) error {
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
		logger.WithError(err).Error("Failed to migrate models")
		return fmt.Errorf("failed to migrate models: %w", err)
	}

	logger.Info("Database migrations completed successfully!")
	return nil
}

// Вспомогательная функция для получения переменной окружения с fallback-значением.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
