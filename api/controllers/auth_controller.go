package controllers

import (
	"net/http"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/dto/mappers"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	service *services.AuthService
	Logger  *logging.Logger
}

func NewAuthController(service *services.AuthService, logger *logging.Logger) *AuthController {
	return &AuthController{service: service, Logger: logger}
}

// SignUp регистрирует нового пользователя.
// @Summary Register a new user
// @Description Register a new user by providing a username, email, and password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.AuthInput true "Registration data"
// @Success 201 {object} dto.AuthResponse "User successfully registered"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/signup [post]
func (c *AuthController) SignUp(ctx *gin.Context) {
	var input dto.AuthInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to bind JSON in SignUp", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.IP = ctx.ClientIP()
	input.UserAgent = ctx.Request.UserAgent()

	user, tokens, err := c.service.SignUp(input)
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to register user in service", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := mappers.MapToAuthResponse(user, tokens.AccessToken, tokens.RefreshToken)
	ctx.JSON(http.StatusCreated, response)
}

// Login аутентифицирует пользователя.
// @Summary Authenticate a user
// @Description Authenticate a user using their email and password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.AuthInput true "Login credentials"
// @Success 200 {object} dto.AuthResponse "User successfully authenticated"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /users/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var input dto.AuthInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to bind JSON in Login", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.IP = ctx.ClientIP()
	input.UserAgent = ctx.Request.UserAgent()

	user, tokens, err := c.service.Login(input)
	if err != nil {
		c.Logger.Log(logrus.WarnLevel, "Authentication failed", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	response := mappers.MapToAuthResponse(user, tokens.AccessToken, tokens.RefreshToken)
	ctx.JSON(http.StatusOK, response)
}

// RefreshToken обновляет access token с использованием refresh token.
// @Summary Refresh tokens
// @Description Refresh the access token by providing the refresh token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.RefreshTokenInput true "Refresh token"
// @Success 200 {object} dto.RefreshTokenResponse "Tokens successfully refreshed"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 401 {object} map[string]string "Invalid or expired refresh token"
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var input dto.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to bind JSON in RefreshToken", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := c.service.RefreshToken(input.RefreshToken)
	if err != nil {
		c.Logger.Log(logrus.WarnLevel, "Failed to refresh token in service", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response := mappers.MapToRefreshTokenResponse(tokens.AccessToken, tokens.RefreshToken)
	ctx.JSON(http.StatusOK, response)
}

// Logout отзывается refresh token.
// @Summary Logout user
// @Description Revoke the refresh token to log out the user.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.LogoutInput true "Refresh token"
// @Success 200 {object} map[string]string "Successfully logged out"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 401 {object} map[string]string "Invalid or not found refresh token"
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	var input dto.LogoutInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to bind JSON in Logout", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Logout(input.RefreshToken); err != nil {
		c.Logger.Log(logrus.WarnLevel, "Failed to logout user", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
