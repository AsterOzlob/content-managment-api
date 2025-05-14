package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// DBConfig содержит настройки для подключения к базе данных.
type DBConfig struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-default:"password"`
	Name     string `env:"DB_NAME" env-default:"content_management"`
	SSLMode  string `env:"DB_SSL_MODE" env-default:"disable"`
}

// LoadDBConfig загружает конфигурацию базы данных из переменных окружения.
func LoadDBConfig() (*DBConfig, error) {
	var cfg DBConfig

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read DB config from environment: %w", err)
	}

	return &cfg, nil
}
