package controllers

import (
	"net/http"
	"strconv"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/dto/mappers"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UserController provides methods for managing users via the HTTP API.
type UserController struct {
	service *services.UserService
	Logger  logger.Logger
}

// NewUserController creates a new instance of UserController.
func NewUserController(service *services.UserService, logger logger.Logger) *UserController {
	return &UserController{
		service: service,
		Logger:  logger,
	}
}

// @Summary Get all users
// @Description Get a list of all users in the system.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.UserResponse "List of users"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	c.Logger.Info("Fetching all users in controller")

	users, err := c.service.GetAllUsers()
	if err != nil {
		c.Logger.WithError(err).Error("Failed to fetch all users in service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, mappers.MapToUserListResponse(users))
}

// @Summary Get user by ID
// @Description Get a user by their unique identifier.
// @Tags Users
// @Produce json
// @Param id path uint true "User ID"
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse "User details"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Logger.WithError(err).Error("Invalid user ID in GetUserByID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	c.Logger.WithField("user_id", id).Info("Fetching user by ID")

	user, err := c.service.GetUserByID(uint(id))
	if err != nil {
		c.Logger.WithError(err).Error("Failed to fetch user by ID in service")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, mappers.MapToUserResponse(user))
}

// @Summary Delete a user
// @Description Delete a user by their unique identifier.
// @Tags Users
// @Produce json
// @Param id path uint true "User ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Logger.WithError(err).Error("Invalid user ID in DeleteUser")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	c.Logger.WithField("target_user_id", id).Info("Deleting user")

	if err := c.service.DeleteUser(uint(id)); err != nil {
		c.Logger.WithError(err).Error("Failed to delete user in service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// @Summary Assign role to a user
// @Description Assigns a role to a user by their unique identifier.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path uint true "User ID"
// @Param role body dto.UserRoleAssignmentInput true "Role name to assign"
// @Security BearerAuth
// @Success 200 {object} map[string]string "Role assigned"
// @Failure 400 {object} map[string]string "Invalid input data or user ID"
// @Failure 500 {object} map[string]string "Failed to assign role"
// @Router /users/{id}/role [patch]
func (c *UserController) AssignRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Logger.WithError(err).Error("Invalid target user ID in AssignRole")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var input dto.UserRoleAssignmentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.WithError(err).Error("Failed to bind JSON in AssignRole")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Logger.WithFields(logrus.Fields{
		"user_id":  id,
		"new_role": input.RoleName,
	}).Info("Assigning role to user")

	if err := c.service.AssignRole(uint(id), input.RoleName); err != nil {
		c.Logger.WithError(err).Error("Failed to assign role in service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "role assigned"})
}
