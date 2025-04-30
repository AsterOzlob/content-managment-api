package dto

// RoleCreateDTO представляет данные для создания роли.
type RoleCreateDTO struct {
	Name        string `json:"name" binding:"required" example:"admin"`  // Название роли.
	Description string `json:"description" example:"Administrator role"` // Описание роли.
}

// RoleUpdateDTO представляет данные для обновления роли.
type RoleUpdateDTO struct {
	Name        string `json:"name" binding:"required" example:"editor"` // Новое название роли.
	Description string `json:"description" example:"Editor role"`        // Новое описание роли.
}

// RoleResponseDTO представляет данные роли для ответа клиенту.
type RoleResponseDTO struct {
	ID          uint   `json:"id" example:"1"`                            // Уникальный идентификатор роли.
	Name        string `json:"name" example:"admin"`                      // Название роли.
	Description string `json:"description" example:"Administrator role"`  // Описание роли.
	CreatedAt   string `json:"created_at" example:"2023-01-01T12:00:00Z"` // Дата создания записи.
	UpdatedAt   string `json:"updated_at" example:"2023-01-01T12:00:00Z"` // Дата последнего обновления записи.
}
