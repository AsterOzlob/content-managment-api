package scheduler

import (
	"time"

	"github.com/AsterOzlob/content_managment_api/internal/database/repositories"
	"github.com/AsterOzlob/content_managment_api/internal/logger"
)

// StartTokenCleanupScheduler запускает планировщик для очистки истекших токенов.
func StartTokenCleanupScheduler(refreshTokenRepo *repositories.RefreshTokenRepository, logger logger.Logger) {
	go func() {
		for {
			time.Sleep(1 * time.Hour) // Запуск каждый час
			logger.Info("Running scheduled cleanup of expired refresh tokens")
			if err := refreshTokenRepo.CleanupExpiredTokens(); err != nil {
				logger.WithError(err).Error("Error during cleanup of expired refresh tokens")
			} else {
				logger.Info("Successfully cleaned up expired refresh tokens")
			}
		}
	}()
}
