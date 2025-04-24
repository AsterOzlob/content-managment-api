package config

import (
	"fmt"
	"os"

	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/sirupsen/logrus"
)

// Config объединяет все конфигурации приложения.
type Config struct {
	DBConfig  *DBConfig  // Конфигурация базы данных.
	JWTConfig *JWTConfig // Конфигурация JWT.
}

// LoadConfig загружает общую конфигурацию приложения.
func LoadConfig(logger *logging.Logger) (*Config, error) {
	// Загрузка конфигураций
	dbConfig, err := LoadDBConfig(logger)
	if err != nil {
		logger.Log(logrus.ErrorLevel, "Failed to load DB config", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to load DB config: %w", err)
	}

	jwtConfig, err := LoadJWTConfig(logger)
	if err != nil {
		logger.Log(logrus.ErrorLevel, "Failed to load JWT config", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to load JWT config: %w", err)
	}

	return &Config{
		DBConfig:  dbConfig,
		JWTConfig: jwtConfig,
	}, nil
}

// Вспомогательная функция для получения переменной окружения с fallback-значением.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
