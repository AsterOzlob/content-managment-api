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
