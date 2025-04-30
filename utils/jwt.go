package utils

import (
	"fmt"
	"time"

	"github.com/AsterOzlob/content_managment_api/config"
	"github.com/golang-jwt/jwt/v5"
)

// RegenerateAccessToken создает новый access token с использованием user_id и роли.
func RegenerateAccessToken(userID uint, role string, cfg *config.JWTConfig) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Duration(cfg.AccessTokenTTL) * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.AccessTokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, nil
}

// GenerateAccessToken создает JWT access token.
func GenerateAccessToken(userID uint, role string, cfg *config.JWTConfig) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Duration(cfg.AccessTokenTTL) * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.AccessTokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken создает JWT refresh token.
func GenerateRefreshToken(userID uint, cfg *config.JWTConfig) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.RefreshTokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateAccessToken проверяет JWT access token.
func ValidateAccessToken(tokenString string, cfg *config.JWTConfig) (*jwt.Token, error) {
	return validateToken(tokenString, cfg.AccessTokenSecret)
}

// ValidateRefreshToken проверяет JWT refresh token.
func ValidateRefreshToken(tokenString string, cfg *config.JWTConfig) (*jwt.Token, error) {
	return validateToken(tokenString, cfg.RefreshTokenSecret)
}

// Вспомогательная функция для валидации токена
func validateToken(tokenString, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return token, nil
}
