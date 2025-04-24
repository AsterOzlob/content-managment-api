package config

import (
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// DBConfig содержит настройки для подключения к базе данных.
type DBConfig struct {
	DBHost     string // Хост базы данных.
	DBPort     string // Порт базы данных.
	DBUser     string // Имя пользователя базы данных.
	DBPassword string // Пароль пользователя базы данных.
	DBName     string // Имя базы данных.
}

// LoadDBConfig загружает конфигурацию базы данных из переменных окружения или .env файла.
func LoadDBConfig(logger *logging.Logger) (*DBConfig, error) {
	// Загрузка .env файла
	if err := godotenv.Load("./.env"); err != nil {
		logger.Log(logrus.WarnLevel, "No .env file found, using environment variables", nil)
	}

	return &DBConfig{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "content_management"),
	}, nil
}
