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

// CommentController предоставляет методы для управления комментариями через HTTP API.
type CommentController struct {
	service *services.CommentService
}

// NewCommentController создаёт новый экземпляр CommentController.
func NewCommentController(service *services.CommentService) *CommentController {
	return &CommentController{service: service}
}

// @Summary Добавить комментарий к статье
// @Description Добавляет новый комментарий к статье по её ID.
// @Tags Комментарии
// @Accept json
// @Produce json
// @Param id path uint true "ID статьи"
// @Param comment body dto.CommentInput true "Данные комментария"
// @Security BearerAuth
// @Success 201 {object} dto.CommentResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles/{id}/comments [post]
func (c *CommentController) AddCommentToArticle(ctx *gin.Context) {
	articleIDStr := ctx.Param("id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrInvalidArticleID})
		return
	}

	var input dto.CommentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUserNotAuthenticated})
		return
	}

	comment, err := c.service.AddCommentToArticle(uint(articleID), input, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrInternalServerError})
		return
	}
	ctx.JSON(http.StatusCreated, mappers.MapToCommentResponse(comment))
}

// @Summary Получить комментарии по ID статьи
// @Description Возвращает все комментарии для указанной статьи, включая вложенные.
// @Tags Комментарии
// @Produce json
// @Param id path uint true "ID статьи"
// @Security BearerAuth
// @Success 200 {array} dto.CommentResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles/{id}/comments [get]
func (c *CommentController) GetCommentsByArticleID(ctx *gin.Context) {
	articleIDStr := ctx.Param("id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrInvalidArticleID})
		return
	}

	comments, err := c.service.GetCommentsByArticleID(uint(articleID))
	if err != nil {
		switch err.Error() {
		case apperrors.ErrArticleNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": apperrors.ErrArticleNotFound})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrInternalServerError})
		}
		return
	}
	ctx.JSON(http.StatusOK, mappers.MapToCommentListResponse(comments))
}

// @Summary Обновить комментарий
// @Description Обновляет существующий комментарий по его ID, если пользователь — владелец, модератор или администратор.
// @Tags Комментарии
// @Accept json
// @Produce json
// @Param id path uint true "ID комментария"
// @Param comment body dto.CommentInput true "Обновлённые данные комментария"
// @Security BearerAuth
// @Success 200 {object} dto.CommentResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles/comments/{id} [put]
func (c *CommentController) UpdateComment(ctx *gin.Context) {
	commentIDStr := ctx.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	var input dto.CommentInput
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

	comment, err := c.service.UpdateComment(uint(commentID), input, userID, userRoles)
	if err != nil {
		switch err.Error() {
		case apperrors.ErrCommentNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": apperrors.ErrCommentNotFound})
		case apperrors.ErrAccessDenied:
			ctx.JSON(http.StatusForbidden, gin.H{"error": apperrors.ErrAccessDenied})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrInternalServerError})
		}
		return
	}
	ctx.JSON(http.StatusOK, mappers.MapToCommentResponse(comment))
}

// @Summary Удалить комментарий
// @Description Удаляет комментарий по его уникальному идентификатору.
// @Tags Комментарии
// @Produce json
// @Param id path uint true "ID комментария"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /comments/{id} [delete]
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	commentIDStr := ctx.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
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

	if err := c.service.DeleteComment(uint(commentID), userID, userRoles); err != nil {
		switch err.Error() {
		case apperrors.ErrCommentNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": apperrors.ErrCommentNotFound})
		case apperrors.ErrAccessDenied:
			ctx.JSON(http.StatusForbidden, gin.H{"error": apperrors.ErrAccessDenied})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrInternalServerError})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}
