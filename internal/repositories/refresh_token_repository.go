package repositories

import (
	"errors"
	"time"

	"github.com/AsterOzlob/content_managment_api/internal/models"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	DB     *gorm.DB
	Logger *logging.Logger
}

func NewRefreshTokenRepository(db *gorm.DB, logger *logging.Logger) *RefreshTokenRepository {
	return &RefreshTokenRepository{DB: db, Logger: logger}
}

// Create создает новый refresh token в базе данных.
func (r *RefreshTokenRepository) Create(token *models.RefreshToken) error {
	r.Logger.Log(logrus.InfoLevel, "Creating new refresh token in database", map[string]interface{}{
		"user_id": token.UserID,
	})
	result := r.DB.Create(token)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to create refresh token in database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}
	return nil
}

func (r *RefreshTokenRepository) GetActiveTokenByUser(userID uint) (*models.RefreshToken, error) {
	var token models.RefreshToken
	result := r.DB.Where("user_id = ? AND expires_at > ? AND revoked = ?", userID, time.Now(), false).First(&token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.Logger.Log(logrus.ErrorLevel, "Failed to fetch active refresh token by user", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}
	return &token, nil
}

// GetByToken находит refresh token по его значению.
func (r *RefreshTokenRepository) GetByToken(tokenString string) (*models.RefreshToken, error) {
	r.Logger.Log(logrus.InfoLevel, "Fetching refresh token by token string", map[string]interface{}{
		"token": tokenString,
	})
	var token models.RefreshToken
	result := r.DB.Where("token = ?", tokenString).First(&token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.Logger.Log(logrus.WarnLevel, "Refresh token not found", nil)
			return nil, nil
		}
		r.Logger.Log(logrus.ErrorLevel, "Failed to fetch refresh token from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return nil, result.Error
	}
	return &token, nil
}

// Update обновляет refresh token в базе данных.
func (r *RefreshTokenRepository) Update(token *models.RefreshToken) error {
	r.Logger.Log(logrus.InfoLevel, "Updating refresh token in database", map[string]interface{}{
		"user_id": token.UserID,
	})
	result := r.DB.Save(token)
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to update refresh token in database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}
	return nil
}

// CleanupExpiredTokens удаляет все истекшие refresh-токены.
func (r *RefreshTokenRepository) CleanupExpiredTokens() error {
	r.Logger.Log(logrus.InfoLevel, "Cleaning up expired refresh tokens", nil)
	result := r.DB.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{})
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to cleanup expired refresh tokens", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}
	r.Logger.Log(logrus.InfoLevel, "Successfully cleaned up expired refresh tokens", map[string]interface{}{
		"deleted_count": result.RowsAffected,
	})
	return nil
}

// Delete удаляет refresh token из базы данных.
func (r *RefreshTokenRepository) Delete(tokenString string) error {
	r.Logger.Log(logrus.InfoLevel, "Deleting refresh token from database", map[string]interface{}{
		"token": tokenString,
	})
	result := r.DB.Where("token = ?", tokenString).Delete(&models.RefreshToken{})
	if result.Error != nil {
		r.Logger.Log(logrus.ErrorLevel, "Failed to delete refresh token from database", map[string]interface{}{
			"error": result.Error.Error(),
		})
		return result.Error
	}
	return nil
}
