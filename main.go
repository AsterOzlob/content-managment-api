package main

import (
	"fmt"

	"github.com/AsterOzlob/content_managment_api/config"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// Инициализация подключения к БД
	dbConn, err := config.InitDB(&cfg)
	if err != nil {
		fmt.Println("Error initializing database connection:", err)
		return
	}

	// Миграция моделей
	err = config.MigrateModels(dbConn)
	if err != nil {
		fmt.Println("Error migrating models:", err)
		return
	}

	fmt.Println("Application started successfully!")
}
