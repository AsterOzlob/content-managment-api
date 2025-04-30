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

type RoleController struct {
	service *services.RoleService
	Logger  *logging.Logger
}

func NewRoleController(service *services.RoleService, logger *logging.Logger) *RoleController {
	return &RoleController{service: service, Logger: logger}
}

// @Summary Создание новой роли
// @Description Создает новую роль в системе.
// @Tags Roles
// @Accept json
// @Produce json
// @Param role body dto.RoleCreateDTO true "Role data"
// @Success 201 {object} dto.RoleResponseDTO "Role created successfully"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /roles [post]
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var input dto.RoleCreateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.Log(logrus.WarnLevel, "Invalid input data for role creation", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Logger.Log(logrus.InfoLevel, "Handling role creation request", map[string]interface{}{
		"role_name": input.Name,
	})
	createdRole, err := c.service.CreateRole(&input)
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to handle role creation request", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, mappers.ToRoleResponseDTO(createdRole))
}

// @Summary Получение всех ролей
// @Description Возвращает список всех ролей в системе.
// @Tags Roles
// @Produce json
// @Success 200 {array} dto.RoleResponseDTO "List of roles"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /roles [get]
func (c *RoleController) GetAllRoles(ctx *gin.Context) {
	c.Logger.Log(logrus.InfoLevel, "Handling request to fetch all roles", nil)

	// Получаем роли из сервиса
	roles, err := c.service.GetAllRoles()
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to handle request to fetch all roles", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем модели в DTO
	var roleDTOs []*dto.RoleResponseDTO
	for _, role := range roles {
		roleDTOs = append(roleDTOs, mappers.ToRoleResponseDTO(role))
	}

	// Отправляем DTO клиенту
	ctx.JSON(http.StatusOK, roleDTOs)
}

// @Summary Получение роли по ID
// @Description Возвращает роль по её уникальному идентификатору.
// @Tags Roles
// @Produce json
// @Param id path uint true "Role ID"
// @Success 200 {object} dto.RoleResponseDTO "Role details"
// @Failure 400 {object} map[string]string "Invalid role ID"
// @Failure 404 {object} map[string]string "Role not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /roles/{id} [get]
func (c *RoleController) GetRoleByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.Logger.Log(logrus.WarnLevel, "Invalid role ID in request", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}
	c.Logger.Log(logrus.InfoLevel, "Handling role fetch request by ID", map[string]interface{}{
		"role_id": id,
	})
	role, err := c.service.GetRoleByID(uint(id))
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to handle role fetch request by ID", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}
	ctx.JSON(http.StatusOK, mappers.ToRoleResponseDTO(role))
}

// @Summary Обновление роли
// @Description Обновляет существующую роль в системе.
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path uint true "Role ID"
// @Param role body dto.RoleUpdateDTO true "Updated role data"
// @Success 200 {object} dto.RoleResponseDTO "Role updated successfully"
// @Failure 400 {object} map[string]string "Invalid input data or role ID"
// @Failure 404 {object} map[string]string "Role not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /roles/{id} [put]
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.Logger.Log(logrus.WarnLevel, "Invalid role ID in request", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}
	var input dto.RoleUpdateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.Log(logrus.WarnLevel, "Invalid input data for role update", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Logger.Log(logrus.InfoLevel, "Handling role update request", map[string]interface{}{
		"role_id": id,
	})
	updatedRole, err := c.service.UpdateRole(uint(id), &input)
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to handle role update request", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mappers.ToRoleResponseDTO(updatedRole))
}

// @Summary Удаление роли
// @Description Удаляет роль по её уникальному идентификатору.
// @Tags Roles
// @Produce json
// @Param id path uint true "Role ID"
// @Success 200 {object} map[string]string "Role deleted successfully"
// @Failure 400 {object} map[string]string "Invalid role ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /roles/{id} [delete]
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.Logger.Log(logrus.WarnLevel, "Invalid role ID in request", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}
	c.Logger.Log(logrus.InfoLevel, "Handling role deletion request", map[string]interface{}{
		"role_id": id,
	})
	if err := c.service.DeleteRole(uint(id)); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to handle role deletion request", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
