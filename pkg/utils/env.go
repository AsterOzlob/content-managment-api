package utils

import "os"

// Вспомогательная функция для получения переменной окружения с fallback-значением.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
