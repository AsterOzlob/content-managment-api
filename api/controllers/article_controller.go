package controllers

import (
	"net/http"
	"strconv"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/dto/mappers"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	apperrors "github.com/AsterOzlob/content_managment_api/pkg/errors"
	"github.com/AsterOzlob/content_managment_api/pkg/utils"
	"github.com/gin-gonic/gin"
)

// ArticleController предоставляет методы для управления статьями через HTTP API.
type ArticleController struct {
	service *services.ArticleService
}

// NewArticleController создаёт новый экземпляр ArticleController.
func NewArticleController(service *services.ArticleService) *ArticleController {
	return &ArticleController{service: service}
}

// @Summary Создать новую статью
// @Description Создает новую статью с возможностью прикрепления медиафайлов.
// @Tags Статьи
// @Accept json
// @Produce json
// @Param article body dto.ArticleInput true "Данные статьи"
// @Security BearerAuth
// @Success 201 {object} dto.ArticleResponse
// @Failure 400 {object} map[string]string
// @Router /articles [post]
func (c *ArticleController) CreateArticle(ctx *gin.Context) {
	var input dto.ArticleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrInvalidArticleID})
		return
	}
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUserNotAuthenticated})
		return
	}
	article, err := c.service.CreateArticle(input, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrInternalServerError})
		return
	}
	ctx.JSON(http.StatusCreated, mappers.MapToArticleResponse(article))
}

// @Summary Получить все статьи
// @Description Возвращает список всех статей с медиафайлами и комментариями.
// @Tags Статьи
// @Produce json
// @Success 200 {array} dto.ArticleResponse
// @Failure 500 {object} map[string]string
// @Router /articles [get]
func (c *ArticleController) GetAllArticles(ctx *gin.Context) {
	articles, err := c.service.GetAllArticles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrInternalServerError})
		return
	}
	ctx.JSON(http.StatusOK, mappers.MapToArticleListResponse(articles))
}

// @Summary Получить статью по ID
// @Description Возвращает статью по её уникальному идентификатору.
// @Tags Статьи
// @Produce json
// @Param id path uint true "ID статьи"
// @Success 200 {object} dto.ArticleResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles/{id} [get]
func (c *ArticleController) GetArticleByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrInvalidArticleID})
		return
	}
	article, err := c.service.GetArticleByID(uint(id))
	if err != nil {
		switch err.Error() {
		case apperrors.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": apperrors.ErrArticleNotFound})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrInternalServerError})
		}
		return
	}
	ctx.JSON(http.StatusOK, mappers.MapToArticleResponse(article))
}

// @Summary Обновить статью
// @Description Обновляет существующую статью.
// @Tags Статьи
// @Accept json
// @Produce json
// @Param id path uint true "ID статьи"
// @Param article body dto.ArticleInput true "Обновлённые данные статьи"
// @Security BearerAuth
// @Success 200 {object} dto.ArticleResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles/{id} [put]
func (c *ArticleController) UpdateArticle(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrInvalidArticleID})
		return
	}
	var input dto.ArticleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUserNotAuthenticated})
		return
	}
	userRoles, err := utils.GetUserRolesFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUserRolesNotFound})
		return
	}
	article, err := c.service.UpdateArticle(uint(id), input, userID, userRoles)
	if err != nil {
		switch err.Error() {
		case apperrors.ErrAccessDenied:
			ctx.JSON(http.StatusForbidden, gin.H{"error": apperrors.ErrAccessDenied})
		case apperrors.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": apperrors.ErrArticleNotFound})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrInternalServerError})
		}
		return
	}
	ctx.JSON(http.StatusOK, mappers.MapToArticleResponse(article))
}

// @Summary Удалить статью
// @Description Удаляет статью по её уникальному идентификатору.
// @Tags Статьи
// @Produce json
// @Param id path uint true "ID статьи"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles/{id} [delete]
func (c *ArticleController) DeleteArticle(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrInvalidArticleID})
		return
	}
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUserNotAuthenticated})
		return
	}
	userRoles, err := utils.GetUserRolesFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUserRolesNotFound})
		return
	}
	err = c.service.DeleteArticle(uint(id), userID, userRoles)
	if err != nil {
		switch err.Error() {
		case apperrors.ErrAccessDenied:
			ctx.JSON(http.StatusForbidden, gin.H{"error": apperrors.ErrAccessDenied})
		case apperrors.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": apperrors.ErrArticleNotFound})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrInternalServerError})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "article deleted successfully"})
}
