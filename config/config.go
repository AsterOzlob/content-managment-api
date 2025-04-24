package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBConfig  *DBConfig
	JWTConfig *JWTConfig
}

// Вспомогательная функция для получения переменной окружения с fallback-значением
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func LoadConfig() (*Config, error) {
	// Загрузка конфигураций
	dbConfig, err := LoadDBConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load DB config: %w", err)
	}

	jwtConfig, err := LoadJWTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load JWT config: %w", err)
	}

	return &Config{
		DBConfig:  dbConfig,
		JWTConfig: jwtConfig,
	}, nil
}
