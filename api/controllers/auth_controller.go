package controllers

import (
	"net/http"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	dtomappers "github.com/AsterOzlob/content_managment_api/internal/dto/mappers"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	"github.com/gin-gonic/gin"
)

// AuthController предоставляет методы для аутентификации.
type AuthController struct {
	service *services.AuthService
}

// NewAuthController создаёт новый экземпляр AuthController.
func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя с указанием имени, email и пароля.
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param input body dto.AuthInput true "Данные регистрации"
// @Success 201 {object} dto.AuthResponse "Пользователь успешно зарегистрирован"
// @Failure 400 {object} map[string]string "Неверные входные данные"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /users/signup [post]
func (c *AuthController) SignUp(ctx *gin.Context) {
	var input dto.AuthInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.IP = ctx.ClientIP()
	input.UserAgent = ctx.Request.UserAgent()

	user, tokens, err := c.service.SignUp(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := dtomappers.MapToAuthResponse(user, tokens.AccessToken, tokens.RefreshToken)
	ctx.JSON(http.StatusCreated, response)
}

// @Summary Аутентификация пользователя
// @Description Аутентифицирует пользователя по email и паролю.
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param input body dto.AuthInput true "Учётные данные"
// @Success 200 {object} dto.AuthResponse "Аутентификация успешна"
// @Failure 401 {object} map[string]string "Неверные учётные данные"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /users/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var input dto.AuthInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.IP = ctx.ClientIP()
	input.UserAgent = ctx.Request.UserAgent()

	user, tokens, err := c.service.Login(input)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "неверные учётные данные"})
		return
	}
	response := dtomappers.MapToAuthResponse(user, tokens.AccessToken, tokens.RefreshToken)
	ctx.JSON(http.StatusOK, response)
}

// @Summary Обновление токена
// @Description Обновляет access токен с помощью refresh токена.
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param input body dto.RefreshTokenInput true "Refresh токен"
// @Success 200 {object} dto.RefreshTokenResponse "Токены успешно обновлены"
// @Failure 400 {object} map[string]string "Ошибка валидации JSON"
// @Failure 401 {object} map[string]string "Неверный или истёкший refresh токен"
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var input dto.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := c.service.RefreshToken(input.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный или истекший refresh токен"})
		return
	}
	response := dtomappers.MapToRefreshTokenResponse(tokens.AccessToken, tokens.RefreshToken)
	ctx.JSON(http.StatusOK, response)
}

// @Summary Отзыв refresh токена
// @Description Отзывает refresh токен для текущей сессии.
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param input body dto.LogoutInput true "Входные данные для выхода"
// @Success 200 {object} map[string]string "Выход выполнен успешно"
// @Failure 400 {object} map[string]string "Неверный ввод данных"
// @Failure 401 {object} map[string]string "Неверный или не найденный refresh токен"
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	var input dto.LogoutInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Logout(input.RefreshToken); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный или не найденный refresh токен"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "успешно вышли из системы"})
}
