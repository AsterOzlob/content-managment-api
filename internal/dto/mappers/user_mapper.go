package mappers

import (
	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
)

// MapToUserResponse преобразует модель User в DTO UserResponse (без токенов).
func MapToUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// MapToUserListResponse преобразует список пользователей в список DTO.
func MapToUserListResponse(users []*models.User) []dto.UserResponse {
	dtoUsers := make([]dto.UserResponse, 0, len(users))

	for _, user := range users {
		dtoUser := MapToUserResponse(user)
		dtoUsers = append(dtoUsers, *dtoUser)
	}

	return dtoUsers
}
