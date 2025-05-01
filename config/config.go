package config

import (
	"fmt"
	"os"

	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
)

// Config объединяет все конфигурации приложения.
type Config struct {
	DBConfig    *DBConfig
	JWTConfig   *JWTConfig
	MediaConfig *MediaConfig
}

// LoadConfig загружает общую конфигурацию приложения.
func LoadConfig(logger logger.Logger) (*Config, error) {
	dbConfig, err := LoadDBConfig(logger)
	if err != nil {
		logger.WithError(err).Error("Failed to load DB config")
		return nil, fmt.Errorf("failed to load DB config: %w", err)
	}

	jwtConfig, err := LoadJWTConfig(logger)
	if err != nil {
		logger.WithError(err).Error("Failed to load JWT config")
		return nil, fmt.Errorf("failed to load JWT config: %w", err)
	}

	mediaConfig, err := LoadMediaConfig(logger)
	if err != nil {
		logger.WithError(err).Error("Failed to load Media config")
		return nil, fmt.Errorf("failed to load Media config: %w", err)
	}

	return &Config{
		DBConfig:    dbConfig,
		JWTConfig:   jwtConfig,
		MediaConfig: mediaConfig,
	}, nil
}

// Вспомогательная функция для получения переменной окружения с fallback-значением.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
