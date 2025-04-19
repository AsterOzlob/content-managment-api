package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/AsterOzlob/content_managment_api/internal/models"
)

var DB *gorm.DB

func InitDB(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBName,
		config.DBSSLMode,
		config.DBPassword,
	)

	// Подключение к БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверка соединения
	sqlDB, _ := db.DB()
	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func MigrateModels(db *gorm.DB) error {
	models := []interface{}{
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.RefreshToken{},
		&models.Content{},
		&models.Media{},
		&models.Comment{},
	}

	// Миграция моделей
	err := db.AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("failed to migrate models: %w", err)
	}

	return nil
}
