package controllers

import (
	"net/http"
	"strconv"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/dto/mappers"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	"github.com/AsterOzlob/content_managment_api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ArticleController предоставляет методы для управления статьями через HTTP API.
type ArticleController struct {
	service *services.ArticleService
	Logger  logger.Logger // Теперь используется интерфейс, а не указатель
}

// NewArticleController создаёт новый экземпляр ArticleController.
func NewArticleController(service *services.ArticleService, logger logger.Logger) *ArticleController {
	return &ArticleController{service: service, Logger: logger}
}

// @Summary Create a new article
// @Description Create a new article with optional media attachments.
// @Tags Articles
// @Accept json
// @Produce json
// @Param article body dto.ArticleInput true "Article Data"
// @Security BearerAuth
// @Success 201 {object} dto.ArticleResponse
// @Failure 400 {object} map[string]string
// @Router /articles [post]
func (c *ArticleController) CreateArticle(ctx *gin.Context) {
	var input dto.ArticleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.WithError(err).Error("Failed to bind JSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	c.Logger.WithFields(logrus.Fields{
		"author_id": userID,
		"title":     input.Title,
	}).Info("Creating new article")

	article, err := c.service.CreateArticle(input, userID)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to create article")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, mappers.MapToArticleResponse(article))
}

// @Summary Get all articles
// @Description Get a list of all articles with media and comments.
// @Tags Articles
// @Produce json
// @Success 200 {array} dto.ArticleResponse
// @Failure 500 {object} map[string]string
// @Router /articles [get]
func (c *ArticleController) GetAllArticles(ctx *gin.Context) {
	c.Logger.Info("Fetching all articles")

	articles, err := c.service.GetAllArticles()
	if err != nil {
		c.Logger.WithError(err).Error("Failed to fetch all articles")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, mappers.MapToArticleListResponse(articles))
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
func (c *ArticleController) GetArticleByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Logger.WithField("id", idStr).Error("Invalid article ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	c.Logger.WithField("article_id", id).Info("Fetching article by ID")

	article, err := c.service.GetArticleByID(uint(id))
	if err != nil {
		c.Logger.WithError(err).Error("Failed to fetch article by ID")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем модель в DTO
	ctx.JSON(http.StatusOK, mappers.MapToArticleResponse(article))
}

// @Summary Update a comment
// @Description Update an existing comment by ID if user is owner, moderator or admin.
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path uint true "Comment ID"
// @Param comment body dto.CommentInput true "Updated Comment Data"
// @Security BearerAuth
// @Success 200 {object} dto.CommentResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles/comments/{id} [put]
func (c *ArticleController) UpdateArticle(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Logger.WithField("id", idStr).Error("Invalid article ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	var input dto.ArticleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.Logger.WithError(err).Error("Failed to bind JSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userRoles, err := utils.GetUserRolesFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user roles not found"})
		return
	}

	article, err := c.service.UpdateArticle(uint(id), input, userID, userRoles)
	if err != nil {
		if err.Error() == "access denied: you are not the owner or don't have required role" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else if err.Error() == "article not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.Logger.WithError(err).Error("Failed to update article")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, mappers.MapToArticleResponse(article))
}

// @Summary Delete article
// @Description Delete an article by ID.
// @Tags Articles
// @Produce json
// @Param id path uint true "Article ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles/{id} [delete]
func (c *ArticleController) DeleteArticle(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Logger.WithField("id", idStr).Error("Invalid article ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userRoles, err := utils.GetUserRolesFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user roles not found"})
		return
	}

	err = c.service.DeleteArticle(uint(id), userID, userRoles)
	if err != nil {
		if err.Error() == "access denied: you are not the owner or don't have required role" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else if err.Error() == "article not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.Logger.WithError(err).Error("Failed to delete article")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "article deleted successfully"})
}
