package migrations

import (
	"fmt"

	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"gorm.io/gorm"

	"github.com/AsterOzlob/content_managment_api/internal/database/models"
)

// MigrateModels выполняет миграцию моделей.
func MigrateModels(db *gorm.DB, logger logger.Logger) error {
	models := []interface{}{
		&models.User{},
		&models.Role{},
		&models.RefreshToken{},
		&models.Article{},
		&models.Media{},
		&models.Comment{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		logger.WithError(err).Error("Failed to migrate models")
		return fmt.Errorf("failed to migrate models: %w", err)
	}

	return nil
}
