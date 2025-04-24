package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/AsterOzlob/content_managment_api/internal/models"
)

// InitDB инициализирует подключение к базе данных.
func InitDB(dbConfig *DBConfig) (*gorm.DB, error) {
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
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверка соединения
	sqlDB, _ := db.DB()
	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Database connection established successfully!")
	return db, nil
}

// MigrateModels выполняет миграцию моделей.
func MigrateModels(db *gorm.DB) error {
	models := []interface{}{
		&models.User{},
		&models.Role{},
		&models.RefreshToken{},
		&models.Article{},
		&models.Media{},
		&models.Comment{},
	}

	// Миграция моделей
	err := db.AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("failed to migrate models: %w", err)
	}

	fmt.Println("Database migrations completed successfully!")
	return nil
}
