package config

import (
	"strconv"

	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/joho/godotenv"
)

// JWTConfig содержит настройки для работы с JWT-токенами.
type JWTConfig struct {
	AccessTokenSecret  string // Секрет для генерации access-токенов.
	RefreshTokenSecret string // Секрет для генерации refresh-токенов.
	AccessTokenTTL     int    // Время жизни access-токена в минутах.
	RefreshTokenTTL    int    // Время жизни refresh-токена в минутах.
}

// LoadJWTConfig загружает конфигурацию JWT из переменных окружения или .env файла.
func LoadJWTConfig(logger logger.Logger) (*JWTConfig, error) {
	// Загрузка переменных окружения из .env
	if err := godotenv.Load("./.env"); err != nil {
		logger.Warn("No .env file found, using environment variables")
	}

	// Чтение переменных окружения
	accessTokenTTL, err := strconv.Atoi(getEnv("JWT_ACCESS_TOKEN_TTL", "15")) // Значение по умолчанию: 15 минут
	if err != nil {
		logger.WithError(err).Warn("Invalid JWT_ACCESS_TOKEN_TTL value, using default 15 minutes")
		accessTokenTTL = 15
	}

	refreshTokenTTL, err := strconv.Atoi(getEnv("JWT_REFRESH_TOKEN_TTL", "4320")) // Значение по умолчанию: 72 часа
	if err != nil {
		logger.WithError(err).Warn("Invalid JWT_REFRESH_TOKEN_TTL value, using default 4320 minutes (72 hours)")
		refreshTokenTTL = 4320
	}

	return &JWTConfig{
		AccessTokenSecret:  getEnv("JWT_ACCESS_TOKEN_SECRET", "default_access_secret"),
		RefreshTokenSecret: getEnv("JWT_REFRESH_TOKEN_SECRET", "default_refresh_secret"),
		AccessTokenTTL:     accessTokenTTL,
		RefreshTokenTTL:    refreshTokenTTL,
	}, nil
}
