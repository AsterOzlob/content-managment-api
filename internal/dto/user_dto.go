package dto

import "time"

// UserResponse используется для ответа с данными пользователя.
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRegistrationInput используется для входных данных при регистрации.
type UserRegistrationInput struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginInput используется для входных данных при авторизации.
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserRoleAssignmentInput используется для назначения роли пользователю.
type UserRoleAssignmentInput struct {
	RoleName string `json:"role_name" binding:"required"`
}
