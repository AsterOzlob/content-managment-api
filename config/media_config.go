package config

import (
	"fmt"
	"strings"

	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
)

// MediaConfig содержит настройки для работы с медиафайлами.
type MediaConfig struct {
	StoragePath  string   // Путь для хранения файлов
	AllowedTypes []string // Разрешенные типы файлов
	MaxSize      int64    // Максимальный размер файла (в байтах)
}

// LoadMediaConfig загружает конфигурацию для медиафайлов.
func LoadMediaConfig(logger logger.Logger) (*MediaConfig, error) {
	storagePath := getEnv("MEDIA_STORAGE_PATH", "./uploads")
	allowedTypesStr := getEnv("MEDIA_ALLOWED_TYPES", "image/jpeg,image/png,application/pdf")
	maxSizeStr := getEnv("MEDIA_MAX_SIZE", "5242880") // Default: 5 MB (5_242_880 байт)

	// Парсинг MaxSize
	var maxSizeInt int64
	_, err := fmt.Sscanf(maxSizeStr, "%d", &maxSizeInt)
	if err != nil {
		logger.WithError(err).Warn("Invalid MEDIA_MAX_SIZE value, using default 5 MB (5242880 bytes)")
		maxSizeInt = 5242880
	}

	return &MediaConfig{
		StoragePath:  storagePath,
		AllowedTypes: splitString(allowedTypesStr, ","),
		MaxSize:      maxSizeInt,
	}, nil
}

// Вспомогательная функция для разбиения строки по разделителю.
func splitString(input, delimiter string) []string {
	if input == "" {
		return []string{}
	}
	return strings.Split(input, delimiter)
}
