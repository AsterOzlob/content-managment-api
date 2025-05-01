package controllers

import (
	"net/http"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	dtomappers "github.com/AsterOzlob/content_managment_api/internal/dto/mappers"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	"github.com/gin-gonic/gin"
)

// AuthController предоставляет методы для аутентификации.
type AuthController struct {
	service *services.AuthService
	Logger  logger.Logger
}

// NewAuthController создаёт новый экземпляр AuthController.
func NewAuthController(service *services.AuthService, logger logger.Logger) *AuthController {
	return &AuthController{service: service, Logger: logger}
}

// @Summary Register a new user
// @Description Register a new user by providing username, email, and password.
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
		c.Logger.WithError(err).Error("Failed to bind JSON in SignUp")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.IP = ctx.ClientIP()
	input.UserAgent = ctx.Request.UserAgent()

	user, tokens, err := c.service.SignUp(input)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to register user in service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dtomappers.MapToAuthResponse(user, tokens.AccessToken, tokens.RefreshToken)
	ctx.JSON(http.StatusCreated, response)
}

// @Summary Authenticate a user
// @Description Authenticate a user with their email and password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.AuthInput true "Login credentials"
// @Success 200 {object} dto.AuthResponse "Authentication successful"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var input dto.AuthInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.WithError(err).Error("Failed to bind JSON in Login")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.IP = ctx.ClientIP()
	input.UserAgent = ctx.Request.UserAgent()

	user, tokens, err := c.service.Login(input)
	if err != nil {
		c.Logger.WithError(err).Warn("Authentication failed")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	response := dtomappers.MapToAuthResponse(user, tokens.AccessToken, tokens.RefreshToken)
	ctx.JSON(http.StatusOK, response)
}

// @Summary Refresh tokens
// @Description Refresh the access token using the refresh token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.RefreshTokenInput true "Refresh token input"
// @Success 200 {object} dto.RefreshTokenResponse "Tokens refreshed successfully"
// @Failure 400 {object} map[string]string "Failed to bind JSON"
// @Failure 401 {object} map[string]string "Invalid or expired refresh token"
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var input dto.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.WithError(err).Error("Failed to bind JSON in RefreshToken")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Logger.WithField("refresh_token", input.RefreshToken).Info("Refreshing token in controller")

	tokens, err := c.service.RefreshToken(input.RefreshToken)
	if err != nil {
		c.Logger.WithError(err).Warn("Failed to refresh token in controller")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	response := dtomappers.MapToRefreshTokenResponse(tokens.AccessToken, tokens.RefreshToken)
	ctx.JSON(http.StatusOK, response)
}

// @Summary Revoke refresh token
// @Description Revoke refresh token for current session.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.LogoutInput true "Logout input"
// @Success 200 {object} map[string]string "Successfully logged out"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 401 {object} map[string]string "Invalid or not found refresh token"
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	var input dto.LogoutInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.WithError(err).Error("Failed to bind JSON in Logout")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Logger.WithField("refresh_token", input.RefreshToken).Info("Revoking refresh token")

	if err := c.service.Logout(input.RefreshToken); err != nil {
		c.Logger.WithError(err).Warn("Failed to logout user")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or not found refresh token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
