package mappers

import (
	"time"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/models"
)

// ToRoleModel преобразует RoleCreateDTO в модель Role.
func ToRoleModel(dto *dto.RoleCreateDTO) *models.Role {
	return &models.Role{
		Name:        dto.Name,
		Description: dto.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// ToRoleResponseDTO преобразует модель Role в RoleResponseDTO.
func ToRoleResponseDTO(role *models.Role) *dto.RoleResponseDTO {
	return &dto.RoleResponseDTO{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   role.UpdatedAt.Format(time.RFC3339),
	}
}
