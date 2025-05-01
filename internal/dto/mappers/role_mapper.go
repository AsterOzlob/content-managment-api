package mappers

import (
	"time"

	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
)

// MapToRoleResponse преобразует модель Role в RoleResponseDTO.
func MapToRoleResponse(role *models.Role) *dto.RoleResponseDTO {
	return &dto.RoleResponseDTO{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   role.UpdatedAt.Format(time.RFC3339),
	}
}

// MapToRoleListResponse преобразует список моделей Role в список DTO RoleResponseDTO.
func MapToRoleListResponse(roles []*models.Role) []*dto.RoleResponseDTO {
	dtoRoles := make([]*dto.RoleResponseDTO, 0, len(roles))

	for _, role := range roles {
		dtoRoles = append(dtoRoles, MapToRoleResponse(role))
	}

	return dtoRoles
}
