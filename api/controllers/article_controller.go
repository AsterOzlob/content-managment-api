package controllers

import (
	"net/http"
	"strconv"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/dto/mappers"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/AsterOzlob/content_managment_api/internal/services"
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

	c.Logger.WithFields(logrus.Fields{
		"author_id": input.AuthorID,
		"title":     input.Title,
	}).Info("Creating new article")

	// Создаем новую статью
	article, err := c.service.CreateArticle(dto.ArticleInput{
		AuthorID:  input.AuthorID,
		Title:     input.Title,
		Text:      input.Text,
		Published: input.Published,
		MediaIDs:  input.MediaIDs,
	})
	if err != nil {
		c.Logger.WithError(err).Error("Failed to create article")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем модель в DTO
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

// @Summary Get article by ID
// @Description Get specific article by ID with media and comments.
// @Tags Articles
// @Produce json
// @Param id path uint true "Article ID"
// @Success 200 {object} dto.ArticleResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /articles/{id} [get]
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

// @Summary Update article
// @Description Update an existing article with optional media updates.
// @Tags Articles
// @Accept json
// @Produce json
// @Param id path uint true "Article ID"
// @Param article body dto.ArticleInput true "Updated Article Data"
// @Security BearerAuth
// @Success 200 {object} dto.ArticleResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles/{id} [put]
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

	c.Logger.WithFields(logrus.Fields{
		"article_id": id,
		"title":      input.Title,
	}).Info("Updating article")

	article, err := c.service.UpdateArticle(uint(id), input)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to update article")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

	c.Logger.WithField("article_id", id).Info("Deleting article")

	if err := c.service.DeleteArticle(uint(id)); err != nil {
		c.Logger.WithError(err).Error("Failed to delete article")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "article deleted"})
}
