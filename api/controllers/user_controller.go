package controllers

import (
	"net/http"
	"strconv"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/dto/mappers"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	"github.com/gin-gonic/gin"
)

// UserController предоставляет методы для управления пользователями через HTTP API.
type UserController struct {
	service *services.UserService
}

// NewUserController создаёт новый экземпляр UserController.
func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

// @Summary Получить всех пользователей
// @Description Возвращает список всех пользователей в системе.
// @Tags Пользователи
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.UserResponse "Список пользователей"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /users [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.service.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mappers.MapToUserListResponse(users))
}

// @Summary Получить пользователя по ID
// @Description Возвращает пользователя по его уникальному идентификатору.
// @Tags Пользователи
// @Produce json
// @Param id path uint true "ID пользователя"
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse "Детали пользователя"
// @Failure 400 {object} map[string]string "Неверный ID пользователя"
// @Failure 404 {object} map[string]string "Пользователь не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID"})
		return
	}

	user, err := c.service.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mappers.MapToUserResponse(user))
}

// @Summary Удалить пользователя
// @Description Удаляет пользователя по его уникальному идентификатору.
// @Tags Пользователи
// @Produce json
// @Param id path uint true "ID пользователя"
// @Security BearerAuth
// @Success 200 {object} map[string]string "Пользователь успешно удален"
// @Failure 400 {object} map[string]string "Неверный ID пользователя"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID"})
		return
	}

	if err := c.service.DeleteUser(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "пользователь успешно удален"})
}

// @Summary Назначить роль пользователю
// @Description Назначает роль пользователю по его уникальному идентификатору.
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param id path uint true "ID пользователя"
// @Param role body dto.UserRoleAssignmentInput true "Название роли"
// @Security BearerAuth
// @Success 200 {object} map[string]string "Роль назначена"
// @Failure 400 {object} map[string]string "Неверный ввод или ID пользователя"
// @Failure 500 {object} map[string]string "Не удалось назначить роль"
// @Router /users/{id}/role [patch]
func (c *UserController) AssignRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID"})
		return
	}

	var input dto.UserRoleAssignmentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.AssignRole(uint(id), input.RoleName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "роль успешно назначена"})
}
