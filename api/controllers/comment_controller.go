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

// CommentController предоставляет методы для управления комментариями через HTTP API.
type CommentController struct {
	service *services.CommentService
	Logger  *logging.Logger
}

// NewCommentController создает новый экземпляр CommentController.
func NewCommentController(service *services.CommentService, logger *logging.Logger) *CommentController {
	return &CommentController{service: service, Logger: logger}
}

// @Summary Add a comment to an article
// @Description Add a new comment to an article by its ID.
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path uint true "Article ID"
// @Param comment body dto.CommentInput true "Comment Data"
// @Success 201 {object} dto.CommentResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles/{id}/comments [post]
func (c *CommentController) AddCommentToArticle(ctx *gin.Context) {
	articleIDStr := ctx.Param("id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 64)
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Invalid article ID", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	var input dto.CommentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to bind JSON", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := c.service.AddCommentToArticle(uint(articleID), input)
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to add comment to article", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, mappers.MapToCommentResponse(comment))
}

// @Summary Get comments by article ID
// @Description Get all comments for an article, including nested comments.
// @Tags Comments
// @Produce json
// @Param id path uint true "Article ID"
// @Success 200 {array} dto.CommentResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles/{id}/comments [get]
func (c *CommentController) GetCommentsByArticleID(ctx *gin.Context) {
	articleIDStr := ctx.Param("id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 64)
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Invalid article ID", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	c.Logger.Log(logrus.InfoLevel, "Fetching comments by article ID in controller", map[string]interface{}{
		"article_id": articleID,
	})

	comments, err := c.service.GetCommentsByArticleID(uint(articleID))
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to fetch comments by article ID", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем модели в DTO
	var dtoComments []dto.CommentResponse
	for _, comment := range comments {
		dtoComments = append(dtoComments, mappers.MapToCommentResponse(&comment))
	}

	ctx.JSON(http.StatusOK, dtoComments)
}

// @Summary Delete a comment
// @Description Delete a comment by its ID.
// @Tags Comments
// @Produce json
// @Param id path uint true "Comment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /comments/{id} [delete]
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	commentIDStr := ctx.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Invalid comment ID", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	if err := c.service.DeleteComment(uint(commentID)); err != nil {
		c.Logger.Log(logrus.ErrorLevel, "Failed to delete comment", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}
