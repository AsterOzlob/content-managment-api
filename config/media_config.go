package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// MediaConfig содержит настройки для работы с медиафайлами.
type MediaConfig struct {
	StoragePath  string `env:"MEDIA_STORAGE_PATH" env-default:"./uploads"`
	AllowedTypes string `env:"MEDIA_ALLOWED_TYPES" env-default:"image/jpeg,image/png,application/pdf"`
	MaxSize      int64  `env:"MEDIA_MAX_SIZE" env-default:"5242880"` // in bytes
}

// LoadMediaConfig загружает конфигурацию для медиафайлов из переменных окружения.
func LoadMediaConfig() (*MediaConfig, error) {
	var cfg MediaConfig

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read Media config from environment: %w", err)
	}

	return &cfg, nil
}
