package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// JWTConfig содержит настройки для работы с JWT-токенами.
type JWTConfig struct {
	AccessTokenSecret  string `env:"JWT_ACCESS_TOKEN_SECRET" env-default:"default_access_secret"`
	RefreshTokenSecret string `env:"JWT_REFRESH_TOKEN_SECRET" env-default:"default_refresh_secret"`
	AccessTokenTTL     int    `env:"JWT_ACCESS_TOKEN_TTL" env-default:"15"`    // in minutes
	RefreshTokenTTL    int    `env:"JWT_REFRESH_TOKEN_TTL" env-default:"4320"` // in minutes
}

// LoadJWTConfig загружает конфигурацию JWT из переменных окружения.
func LoadJWTConfig() (*JWTConfig, error) {
	var cfg JWTConfig

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWT config from environment: %w", err)
	}

	return &cfg, nil
}
