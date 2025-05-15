package dto

// RoleCreateDTO представляет данные для создания роли.
type RoleCreateDTO struct {
	Name        string `json:"name" binding:"required"` // Название роли.
	Description string `json:"description"`             // Описание роли.
}

// RoleUpdateDTO представляет данные для обновления роли.
type RoleUpdateDTO struct {
	Name        string `json:"name" binding:"required"` // Новое название роли.
	Description string `json:"description"`             // Новое описание роли.
}

// RoleResponseDTO представляет данные роли для ответа клиенту.
type RoleResponseDTO struct {
	ID          uint   `json:"id"`          // Уникальный идентификатор роли.
	Name        string `json:"name"`        // Название роли.
	Description string `json:"description"` // Описание роли.
	CreatedAt   string `json:"created_at"`  // Дата создания записи.
	UpdatedAt   string `json:"updated_at"`  // Дата последнего обновления записи.
}
