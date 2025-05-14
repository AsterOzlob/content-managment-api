package controllers

import (
	"net/http"
	"strconv"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/dto/mappers"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	"github.com/gin-gonic/gin"
)

// RoleController предоставляет CRUD операции для ролей.
type RoleController struct {
	service *services.RoleService
}

// NewRoleController создаёт новый экземпляр RoleController.
func NewRoleController(service *services.RoleService) *RoleController {
	return &RoleController{service: service}
}

// @Summary Создание новой роли
// @Description Создает новую роль в системе.
// @Tags Роли
// @Accept json
// @Produce json
// @Param role body dto.RoleCreateDTO true "Данные роли"
// @Success 201 {object} dto.RoleResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /roles [post]
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var input dto.RoleCreateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdRole, err := c.service.CreateRole(&input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, mappers.MapToRoleResponse(createdRole))
}

// @Summary Получение всех ролей
// @Description Возвращает список всех ролей в системе.
// @Tags Роли
// @Produce json
// @Success 200 {array} dto.RoleResponseDTO
// @Failure 500 {object} map[string]string
// @Router /roles [get]
func (c *RoleController) GetAllRoles(ctx *gin.Context) {
	roles, err := c.service.GetAllRoles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mappers.MapToRoleListResponse(roles))
}

// @Summary Получение роли по ID
// @Description Возвращает роль по её уникальному идентификатору.
// @Tags Роли
// @Produce json
// @Param id path uint true "ID роли"
// @Success 200 {object} dto.RoleResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /roles/{id} [get]
func (c *RoleController) GetRoleByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID"})
		return
	}

	role, err := c.service.GetRoleByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mappers.MapToRoleResponse(role))
}

// @Summary Обновление роли
// @Description Обновляет существующую роль в системе.
// @Tags Роли
// @Accept json
// @Produce json
// @Param id path uint true "ID роли"
// @Param role body dto.RoleUpdateDTO true "Обновлённые данные роли"
// @Success 200 {object} dto.RoleResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /roles/{id} [put]
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID"})
		return
	}

	var input dto.RoleUpdateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedRole, err := c.service.UpdateRole(uint(id), &input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mappers.MapToRoleResponse(updatedRole))
}

// @Summary Удаление роли
// @Description Удаляет роль по её уникальному идентификатору.
// @Tags Роли
// @Produce json
// @Param id path uint true "ID роли"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /roles/{id} [delete]
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID"})
		return
	}

	if err := c.service.DeleteRole(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "роль успешно удалена"})
}
