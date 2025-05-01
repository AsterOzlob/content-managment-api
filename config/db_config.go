package config

import (
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/joho/godotenv"
)

// DBConfig содержит настройки для подключения к базе данных.
type DBConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// LoadDBConfig загружает конфигурацию базы данных.
func LoadDBConfig(logger logger.Logger) (*DBConfig, error) {
	if err := godotenv.Load("./.env"); err != nil {
		logger.Warn("No .env file found, using environment variables")
	}

	return &DBConfig{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "content_management"),
	}, nil
}
