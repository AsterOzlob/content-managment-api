package config

import (
	"fmt"

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
	dbConfig, err := LoadDBConfig()
	if err != nil {
		logger.WithError(err).Error("Failed to load DB config")
		return nil, fmt.Errorf("failed to load DB config: %w", err)
	}

	jwtConfig, err := LoadJWTConfig()
	if err != nil {
		logger.WithError(err).Error("Failed to load JWT config")
		return nil, fmt.Errorf("failed to load JWT config: %w", err)
	}

	mediaConfig, err := LoadMediaConfig()
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
