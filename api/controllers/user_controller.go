package controllers

import (
	"net/http"
	"strconv"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	"github.com/gin-gonic/gin"
)

// UserController предоставляет методы для управления пользователями через HTTP API.
type UserController struct {
	service *services.UserService // service - экземпляр UserService для выполнения бизнес-логики.
}

// NewUserController создает новый экземпляр UserController.
func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
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
// @x-example {"username":"testuser","email":"test@example.com","password":"SecurePassword123"}
func (c *UserController) SignUp(ctx *gin.Context) {
	var input dto.UserRegistrationInput
	// Привязываем входные данные из тела запроса к структуре UserRegistrationInput.
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вызываем метод сервиса для регистрации пользователя.
	user, err := c.service.SignUp(input.Username, input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем успешный ответ с данными созданного пользователя.
	ctx.JSON(http.StatusCreated, user)
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
// @x-example {"email":"test@example.com","password":"SecurePassword123"}
func (c *UserController) Login(ctx *gin.Context) {
	var input dto.LoginInput
	// Привязываем входные данные из тела запроса к структуре LoginInput.
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вызываем метод сервиса для аутентификации пользователя.
	user, err := c.service.Login(input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем успешный ответ с данными пользователя.
	ctx.JSON(http.StatusOK, user)
}

// DeleteUser удаляет пользователя по ID.
// Только сам пользователь или администратор могут выполнять это действие.
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
	// Получаем ID целевого пользователя из параметров пути.
	targetUserIDStr := ctx.Param("id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 64)
	if err != nil {
		// Если ID некорректен, возвращаем ошибку 400 Bad Request.
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Вызываем метод сервиса для удаления пользователя.
	if err := c.service.DeleteUser(uint(targetUserID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем успешный ответ об удалении пользователя.
	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// GetAllUsers возвращает список всех пользователей.
// Только администратор может выполнять это действие.
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
	// Вызываем метод сервиса для получения всех пользователей.
	users, err := c.service.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем успешный ответ со списком пользователей.
	ctx.JSON(http.StatusOK, users)
}

// GetUserByID возвращает пользователя по ID.
// Только администратор может выполнять это действие.
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path uint true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	// Получаем ID пользователя из параметров пути.
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// Если ID некорректен, возвращаем ошибку 400 Bad Request.
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Вызываем метод сервиса для получения пользователя по ID.
	user, err := c.service.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем успешный ответ с данными пользователя.
	ctx.JSON(http.StatusOK, user)
}

// AssignRole назначает роль пользователю.
// Только администратор может выполнять это действие.
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
// @x-example {"role_name":"admin"}
func (c *UserController) AssignRole(ctx *gin.Context) {
	// Получаем ID целевого пользователя из параметров пути.
	targetUserIDStr := ctx.Param("id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 64)
	if err != nil {
		// Если ID некорректен, возвращаем ошибку 400 Bad Request.
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid target user ID"})
		return
	}

	// Привязываем входные данные из тела запроса к структуре UserRoleAssignmentInput.
	var input dto.UserRoleAssignmentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вызываем метод сервиса для назначения роли пользователю.
	if err := c.service.AssignRole(uint(targetUserID), input.RoleName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем успешный ответ о назначении роли.
	ctx.JSON(http.StatusOK, gin.H{"message": "role assigned"})
}
