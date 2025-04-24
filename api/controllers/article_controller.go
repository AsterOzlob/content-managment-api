package controllers

import (
	"net/http"
	"strconv"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/dto/mappers"
	"github.com/AsterOzlob/content_managment_api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ArticleController предоставляет методы для управления статьями через HTTP API.
type ArticleController struct {
	service *services.ArticleService // service - экземпляр ArticleService для выполнения бизнес-логики.
	logger  *logrus.Logger           // logger - экземпляр логгера для ArticleController.
}

// NewArticleController создает новый экземпляр ArticleController.
func NewArticleController(service *services.ArticleService, logger *logrus.Logger) *ArticleController {
	return &ArticleController{service: service, logger: logger}
}

// @Summary Create a new article
// @Description Create a new article with optional media attachments.
// @Tags Articles
// @Accept json
// @Produce json
// @Param article body dto.ArticleInput true "Article Data"
// @Success 201 {object} dto.ArticleResponse
// @Failure 400 {object} map[string]string
// @Router /article [post]
func (c *ArticleController) CreateArticle(ctx *gin.Context) {
	var input dto.ArticleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to bind JSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.logger.WithFields(logrus.Fields{
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
		c.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create article")
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
// @Router /article [get]
func (c *ArticleController) GetAllArticles(ctx *gin.Context) {
	c.logger.Info("Fetching all articles")

	articles, err := c.service.GetAllArticles()
	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to fetch all articles")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем модели в DTO
	var dtoArticles []*dto.ArticleResponse
	for _, article := range articles {
		dtoArticles = append(dtoArticles, mappers.MapToArticleResponse(article))
	}

	ctx.JSON(http.StatusOK, dtoArticles)
}

// @Summary Get article by ID
// @Description Get a specific article by ID with media and comments.
// @Tags Articles
// @Produce json
// @Param id path uint true "Article ID"
// @Success 200 {object} dto.ArticleResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /article/{id} [get]
func (c *ArticleController) GetArticleByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Invalid article ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	c.logger.WithFields(logrus.Fields{
		"article_id": id,
	}).Info("Fetching article by ID")

	article, err := c.service.GetArticleByID(uint(id))
	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to fetch article by ID")
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
// @Success 200 {object} dto.ArticleResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /article/{id} [put]
func (c *ArticleController) UpdateArticle(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Invalid article ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	var input dto.ArticleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to bind JSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.logger.WithFields(logrus.Fields{
		"article_id": id,
		"title":      input.Title,
	}).Info("Updating article")

	article, err := c.service.UpdateArticle(uint(id), input)
	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to update article")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем модель в DTO
	ctx.JSON(http.StatusOK, mappers.MapToArticleResponse(article))
}

// @Summary Delete article
// @Description Delete an article by ID.
// @Tags Articles
// @Produce json
// @Param id path uint true "Article ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /article/{id} [delete]
func (c *ArticleController) DeleteArticle(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Invalid article ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}

	c.logger.WithFields(logrus.Fields{
		"article_id": id,
	}).Info("Deleting article")

	if err := c.service.DeleteArticle(uint(id)); err != nil {
		c.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to delete article")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "article deleted"})
}
