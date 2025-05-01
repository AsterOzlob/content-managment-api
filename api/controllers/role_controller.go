package controllers

import (
	"net/http"
	"strconv"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/dto/mappers"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	"github.com/gin-gonic/gin"
)

// RoleController предоставляет CRUD операции для ролей.
type RoleController struct {
	service *services.RoleService
	Logger  logger.Logger
}

// NewRoleController создаёт новый экземпляр RoleController.
func NewRoleController(service *services.RoleService, logger logger.Logger) *RoleController {
	return &RoleController{service: service, Logger: logger}
}

// @Summary Создание новой роли
// @Description Создает новую роль в системе.
// @Tags Roles
// @Accept json
// @Produce json
// @Param role body dto.RoleCreateDTO true "Role data"
// @Success 201 {object} dto.RoleResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /roles [post]
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var input dto.RoleCreateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.WithError(err).Warn("Invalid input data for role creation")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Logger.WithField("role_name", input.Name).Info("Handling role creation request")

	createdRole, err := c.service.CreateRole(&input)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to handle role creation request")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, mappers.MapToRoleResponse(createdRole))
}

// @Summary Получение всех ролей
// @Description Возвращает список всех ролей в системе.
// @Tags Roles
// @Produce json
// @Success 200 {array} dto.RoleResponseDTO
// @Failure 500 {object} map[string]string
// @Router /roles [get]
func (c *RoleController) GetAllRoles(ctx *gin.Context) {
	c.Logger.Info("Handling request to fetch all roles")

	roles, err := c.service.GetAllRoles()
	if err != nil {
		c.Logger.WithError(err).Error("Failed to handle request to fetch all roles")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, mappers.MapToRoleListResponse(roles))
}

// @Summary Получение роли по ID
// @Description Возвращает роль по её уникальному идентификатору.
// @Tags Roles
// @Produce json
// @Param id path uint true "Role ID"
// @Success 200 {object} dto.RoleResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /roles/{id} [get]
func (c *RoleController) GetRoleByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Logger.WithError(err).WithField("role_id", idStr).Warn("Invalid role ID in request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	c.Logger.WithField("role_id", id).Info("Handling role fetch request by ID")

	role, err := c.service.GetRoleByID(uint(id))
	if err != nil {
		c.Logger.WithError(err).Error("Failed to handle role fetch request by ID")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	ctx.JSON(http.StatusOK, mappers.MapToRoleResponse(role))
}

// @Summary Обновление роли
// @Description Обновляет существующую роль в системе.
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path uint true "Role ID"
// @Param role body dto.RoleUpdateDTO true "Updated role data"
// @Success 200 {object} dto.RoleResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /roles/{id} [put]
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Logger.WithError(err).WithField("role_id", idStr).Warn("Invalid role ID in request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var input dto.RoleUpdateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.WithError(err).Warn("Invalid input data for role update")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Logger.WithField("role_id", id).Info("Handling role update request")

	updatedRole, err := c.service.UpdateRole(uint(id), &input)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to handle role update request")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, mappers.MapToRoleResponse(updatedRole))
}

// @Summary Удаление роли
// @Description Удаляет роль по её уникальному идентификатору.
// @Tags Roles
// @Produce json
// @Param id path uint true "Role ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /roles/{id} [delete]
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Logger.WithError(err).WithField("role_id", idStr).Warn("Invalid role ID in request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	c.Logger.WithField("role_id", id).Info("Handling role deletion request")

	if err := c.service.DeleteRole(uint(id)); err != nil {
		c.Logger.WithError(err).Error("Failed to handle role deletion request")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
