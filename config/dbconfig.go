package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadDBConfig() (*DBConfig, error) {
	// Загрузка .env файла
	err := godotenv.Load("./.env")
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	// Создание конфигурации с значениями по умолчанию
	return &DBConfig{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "content_management"),
	}, nil
}
