package config

import (
	"fmt"
	"strings"

	logging "github.com/AsterOzlob/content_managment_api/logger"
)

// MediaConfig содержит настройки для работы с медиафайлами.
type MediaConfig struct {
	StoragePath  string   // Путь для хранения файлов
	AllowedTypes []string // Разрешенные типы файлов
	MaxSize      int64    // Максимальный размер файла (в байтах)
}

// LoadMediaConfig загружает конфигурацию для медиафайлов.
func LoadMediaConfig(logger *logging.Logger) (*MediaConfig, error) {
	storagePath := getEnv("MEDIA_STORAGE_PATH", "./uploads")
	allowedTypes := getEnv("MEDIA_ALLOWED_TYPES", "image/jpeg,image/png,application/pdf")
	maxSize := getEnv("MEDIA_MAX_SIZE", "5242880") // Default: 5 MB

	// Преобразование maxSize в int64
	var maxSizeInt int64
	fmt.Sscanf(maxSize, "%d", &maxSizeInt)

	return &MediaConfig{
		StoragePath:  storagePath,
		AllowedTypes: splitString(allowedTypes, ","),
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
