package config

import (
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

type JWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     int // В минутах
	RefreshTokenTTL    int // В минутах
}

func LoadJWTConfig() (*JWTConfig, error) {
	// Загрузка переменных окружения из .env
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Чтение переменных окружения
	accessTokenTTL, err := strconv.Atoi(getEnv("JWT_ACCESS_TOKEN_TTL", "15")) // Значение по умолчанию: 15 минут
	if err != nil {
		log.Println("Invalid JWT_ACCESS_TOKEN_TTL value, using default 15 minutes")
		accessTokenTTL = 15
	}

	refreshTokenTTL, err := strconv.Atoi(getEnv("JWT_REFRESH_TOKEN_TTL", "4320")) // Значение по умолчанию: 72 часа
	if err != nil {
		log.Println("Invalid JWT_REFRESH_TOKEN_TTL value, using default 4320 minutes (72 hours)")
		refreshTokenTTL = 4320
	}

	return &JWTConfig{
		AccessTokenSecret:  getEnv("JWT_ACCESS_TOKEN_SECRET", "default_access_secret"),
		RefreshTokenSecret: getEnv("JWT_REFRESH_TOKEN_SECRET", "default_refresh_secret"),
		AccessTokenTTL:     accessTokenTTL,
		RefreshTokenTTL:    refreshTokenTTL,
	}, nil
}
