package controllers

import (
	"net/http"
	"strconv"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/dto/mappers"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UserController предоставляет методы для управления пользователями через HTTP API.
type UserController struct {
	service *services.UserService // service - экземпляр UserService для выполнения бизнес-логики.
	Logger  *logging.Logger       // Logger - экземпляр логгера для UserController.
}

// NewUserController создает новый экземпляр UserController.
func NewUserController(service *services.UserService, logger *logging.Logger) *UserController {
	return &UserController{service: service, Logger: logger}
}

// SignUp регистрирует нового пользователя.
// @Summary SignUp a new user
// @Description SignUp a new user with username, email, and password
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.UserRegistrationInput true "User Data"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Router /users/signup [post]
func (c *UserController) SignUp(ctx *gin.Context) {
	var input dto.UserRegistrationInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to bind JSON in SignUp", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Logger.Log(logrus.InfoLevel, "Registering new user in controller", map[string]interface{}{
		"username": input.Username,
		"email":    input.Email,
	})

	user, err := c.service.SignUp(input.Username, input.Email, input.Password)
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to register user in service", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, mappers.MapToUserResponse(user))
}

// Login аутентифицирует пользователя.
// @Summary Authenticate a user
// @Description Authenticate a user with email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param credentials body dto.LoginInput true "User Credentials"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var input dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to bind JSON in Login", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Logger.Log(logrus.InfoLevel, "Authenticating user in controller", map[string]interface{}{
		"email": input.Email,
	})

	user, err := c.service.Login(input.Email, input.Password)
	if err != nil {
		c.Logger.Log(logrus.WarnLevel, "Authentication failed", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, mappers.MapToUserResponse(user))
}

// DeleteUser удаляет пользователя по ID.
// @Summary Delete a user
// @Description Delete a user by their ID. Only the user themselves or an admin can perform this action.
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path uint true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	targetUserIDStr := ctx.Param("id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 64)
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Invalid user ID in DeleteUser", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	c.Logger.Log(logrus.InfoLevel, "Deleting user in controller", map[string]interface{}{
		"user_id": targetUserID,
	})

	if err := c.service.DeleteUser(uint(targetUserID)); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to delete user in service", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// GetAllUsers возвращает список всех пользователей.
// @Summary Get all users
// @Description Get a list of all users in the system. Only admin can perform this action.
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.UserResponse
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /users [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	c.Logger.Log(logrus.InfoLevel, "Fetching all users in controller", nil)

	users, err := c.service.GetAllUsers()
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to fetch all users in service", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userResponses []*dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, mappers.MapToUserResponse(user))
	}

	ctx.JSON(http.StatusOK, userResponses)
}

// GetUserByID возвращает пользователя по ID.
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path uint true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Invalid user ID in GetUserByID", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	c.Logger.Log(logrus.InfoLevel, "Fetching user by ID in controller", map[string]interface{}{
		"user_id": id,
	})

	user, err := c.service.GetUserByID(uint(id))
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to fetch user by ID in service", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, mappers.MapToUserResponse(user))
}

// AssignRole назначает роль пользователю.
// @Summary Assign a role to a user
// @Description Assign a role to a user by their ID. Only admin can perform this action.
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path uint true "User ID"
// @Param role body dto.UserRoleAssignmentInput true "Role Name"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /users/{id}/role [patch]
func (c *UserController) AssignRole(ctx *gin.Context) {
	targetUserIDStr := ctx.Param("id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 64)
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Invalid target user ID in AssignRole", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid target user ID"})
		return
	}

	var input dto.UserRoleAssignmentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to bind JSON in AssignRole", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Logger.Log(logrus.InfoLevel, "Assigning role to user in controller", map[string]interface{}{
		"user_id": targetUserID,
		"role":    input.RoleName,
	})

	if err := c.service.AssignRole(uint(targetUserID), input.RoleName); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to assign role in service", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "role assigned"})
}
