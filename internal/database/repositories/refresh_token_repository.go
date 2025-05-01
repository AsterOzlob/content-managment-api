package repositories

import (
	"errors"
	"time"

	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"gorm.io/gorm"
)

// RefreshTokenRepository предоставляет методы для работы с refresh токенами.
type RefreshTokenRepository struct {
	DB     *gorm.DB
	Logger logger.Logger
}

// NewRefreshTokenRepository создаёт новый экземпляр RefreshTokenRepository.
func NewRefreshTokenRepository(db *gorm.DB, logger logger.Logger) *RefreshTokenRepository {
	return &RefreshTokenRepository{DB: db, Logger: logger}
}

// Create создаёт новый refresh token в базе данных.
func (r *RefreshTokenRepository) Create(token *models.RefreshToken) error {
	r.Logger.WithField("user_id", token.UserID).Info("Creating new refresh token in database")
	result := r.DB.Create(token)
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to create refresh token in database")
		return result.Error
	}
	return nil
}

// GetActiveTokenByUser получает активный refresh token по ID пользователя.
func (r *RefreshTokenRepository) GetActiveTokenByUser(userID uint) (*models.RefreshToken, error) {
	var token models.RefreshToken
	result := r.DB.Where("user_id = ? AND expires_at > ? AND revoked = ?", userID, time.Now(), false).First(&token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.Logger.WithError(result.Error).Error("Failed to fetch active refresh token by user")
		return nil, result.Error
	}
	return &token, nil
}

// GetByToken находит refresh token по его значению.
func (r *RefreshTokenRepository) GetByToken(tokenString string) (*models.RefreshToken, error) {
	r.Logger.WithField("token", tokenString).Info("Fetching refresh token by token string")

	var token models.RefreshToken
	result := r.DB.Where("token = ?", tokenString).First(&token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.Logger.Warn("Refresh token not found")
			return nil, nil
		}
		r.Logger.WithError(result.Error).Error("Failed to fetch refresh token from database")
		return nil, result.Error
	}
	return &token, nil
}

// Update обновляет refresh token в базе данных.
func (r *RefreshTokenRepository) Update(token *models.RefreshToken) error {
	r.Logger.WithField("user_id", token.UserID).Info("Updating refresh token in database")
	result := r.DB.Save(token)
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to update refresh token in database")
		return result.Error
	}
	return nil
}

// CleanupExpiredTokens удаляет все истекшие refresh-токены.
func (r *RefreshTokenRepository) CleanupExpiredTokens() error {
	r.Logger.Info("Cleaning up expired refresh tokens")
	result := r.DB.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{})
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to cleanup expired refresh tokens")
		return result.Error
	}
	r.Logger.WithField("deleted_count", result.RowsAffected).Info("Successfully cleaned up expired refresh tokens")
	return nil
}

// Delete удаляет refresh token из базы данных.
func (r *RefreshTokenRepository) Delete(tokenString string) error {
	r.Logger.WithField("token", tokenString).Info("Deleting refresh token from database")
	result := r.DB.Where("token = ?", tokenString).Delete(&models.RefreshToken{})
	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("Failed to delete refresh token from database")
		return result.Error
	}
	return nil
}
