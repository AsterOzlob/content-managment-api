package dto

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthInput struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	IP        string `json:"-"`
	UserAgent string `json:"-"`
}

type LogoutInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthResponse представляет ответ после успешной регистрации или входа.
type AuthResponse struct {
	User         UserResponse `json:"user"`          // Информация о пользователе.
	AccessToken  string       `json:"access_token"`  // Access token.
	RefreshToken string       `json:"refresh_token"` // Refresh token.
}

// RefreshTokenResponse представляет ответ после обновления токенов.
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`  // Новый access token.
	RefreshToken string `json:"refresh_token"` // Новый refresh token.
}
