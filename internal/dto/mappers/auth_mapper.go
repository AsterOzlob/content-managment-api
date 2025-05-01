package mappers

import (
	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
)

// MapToAuthResponse преобразует пользователя и токены в DTO для ответа.
func MapToAuthResponse(user *models.User, accessToken, refreshToken string) *dto.AuthResponse {
	return &dto.AuthResponse{
		User:         *MapToUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

// MapToRefreshTokenResponse преобразует новые токены в DTO для ответа.
func MapToRefreshTokenResponse(accessToken, refreshToken string) *dto.RefreshTokenResponse {
	return &dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
