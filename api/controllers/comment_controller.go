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

// CommentController предоставляет методы для управления комментариями через HTTP API.
type CommentController struct {
	service *services.CommentService
	Logger  logger.Logger
}

// NewCommentController создает новый экземпляр CommentController.
func NewCommentController(service *services.CommentService, logger logger.Logger) *CommentController {
	return &CommentController{service: service, Logger: logger}
}

// @Summary Add a comment to an article
// @Description Add a new comment to an article by its ID.
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path uint true "Article ID"
// @Param comment body dto.CommentInput true "Comment Data"
// @Security BearerAuth
// @Success 201 {object} dto.CommentResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles/{id}/comments [post]
func (c *CommentController) AddCommentToArticle(ctx *gin.Context) {
	articleIDStr := ctx.Param("id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 64)
	if err != nil {
		c.Logger.WithError(err).WithField("article_id", articleIDStr).Error("Invalid article ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	var input dto.CommentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.WithError(err).Error("Failed to bind JSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := c.service.AddCommentToArticle(uint(articleID), input)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to add comment to article")
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
		c.Logger.WithError(err).WithField("article_id", articleIDStr).Error("Invalid article ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	c.Logger.WithField("article_id", articleID).Info("Fetching comments by article ID")

	comments, err := c.service.GetCommentsByArticleID(uint(articleID))
	if err != nil {
		c.Logger.WithError(err).Error("Failed to fetch comments by article ID")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, mappers.MapToCommentListResponse(comments))
}

// @Summary Delete a comment
// @Description Delete a comment by its ID.
// @Tags Comments
// @Produce json
// @Param id path uint true "Comment ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /comments/{id} [delete]
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	commentIDStr := ctx.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		c.Logger.WithError(err).WithField("comment_id", commentIDStr).Error("Invalid comment ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	if err := c.service.DeleteComment(uint(commentID)); err != nil {
		c.Logger.WithError(err).Error("Failed to delete comment")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}
