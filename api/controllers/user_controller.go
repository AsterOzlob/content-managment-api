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

type UserController struct {
	service *services.UserService
	Logger  *logging.Logger
}

func NewUserController(service *services.UserService, logger *logging.Logger) *UserController {
	return &UserController{service: service, Logger: logger}
}

// @Summary Получение всех пользователей
// @Description Возвращает список всех пользователей в системе.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.UserResponse "Список пользователей"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
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

// @Summary Получение пользователя по ID
// @Description Возвращает пользователя по его уникальному идентификатору.
// @Tags Users
// @Produce json
// @Param id path uint true "User ID"
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse "Информация о пользователе"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
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

// @Summary Удаление пользователя
// @Description Удаляет пользователя по его уникальному идентификатору.
// @Tags Users
// @Param id path uint true "User ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string "Пользователь удален"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 500 {object} map[string]string "Internal server error"
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

	if err := c.service.DeleteUser(uint(targetUserID)); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to delete user in service", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// @Summary Назначение роли пользователю
// @Description Назначает роль пользователю по его уникальному идентификатору.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path uint true "User ID"
// @Param role body dto.UserRoleAssignmentInput true "Role name to assign"
// @Security BearerAuth
// @Success 200 {object} map[string]string "Роль назначена"
// @Failure 400 {object} map[string]string "Invalid input data or user ID"
// @Failure 500 {object} map[string]string "Internal server error"
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

	if err := c.service.AssignRole(uint(targetUserID), input.RoleName); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to assign role in service", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "role assigned"})
}
