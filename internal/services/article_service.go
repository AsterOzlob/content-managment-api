package services

import (
	"errors"

	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	"github.com/AsterOzlob/content_managment_api/internal/database/repositories"
	"github.com/AsterOzlob/content_managment_api/internal/dto"
	logger "github.com/AsterOzlob/content_managment_api/internal/logger"
	"github.com/sirupsen/logrus"
)

// ArticleService предоставляет методы для управления статьями.
type ArticleService struct {
	repo   *repositories.ArticleRepository
	Logger logger.Logger
}

// NewArticleService создаёт новый экземпляр ArticleService.
func NewArticleService(repo *repositories.ArticleRepository, logger logger.Logger) *ArticleService {
	return &ArticleService{repo: repo, Logger: logger}
}

// CreateArticle создаёт новую статью.
func (s *ArticleService) CreateArticle(input dto.ArticleInput) (*models.Article, error) {
	s.Logger.WithFields(logrus.Fields{
		"author_id": input.AuthorID,
		"title":     input.Title,
	}).Info("Creating article in service")

	// Проверяем, существует ли пользователь с указанным AuthorID
	var user models.User
	if err := s.repo.DB.First(&user, input.AuthorID).Error; err != nil {
		s.Logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"author_id": input.AuthorID,
		}).Error("User not found")
		return nil, errors.New("user not found")
	}

	article := &models.Article{
		AuthorID:  input.AuthorID,
		Title:     input.Title,
		Text:      input.Text,
		Published: input.Published,
	}

	if err := s.repo.Create(article); err != nil {
		s.Logger.WithError(err).Error("Failed to create article in repository")
		return nil, err
	}

	return article, nil
}

// GetAllArticles возвращает список всех статей.
func (s *ArticleService) GetAllArticles() ([]*models.Article, error) {
	s.Logger.Info("Fetching all articles in service")

	articles, err := s.repo.GetAll()
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch all articles from repository")
		return nil, err
	}
	return articles, nil
}

// GetArticleByID возвращает статью по ID.
func (s *ArticleService) GetArticleByID(id uint) (*models.Article, error) {
	s.Logger.WithField("article_id", id).Info("Fetching article by ID in service")

	article, err := s.repo.GetByID(id)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch article by ID from repository")
		return nil, errors.New("article not found")
	}
	return article, nil
}

// UpdateArticle обновляет существующую статью.
func (s *ArticleService) UpdateArticle(id uint, input dto.ArticleInput) (*models.Article, error) {
	s.Logger.WithFields(logrus.Fields{
		"article_id": id,
		"title":      input.Title,
	}).Info("Updating article in service")

	article, err := s.repo.GetByID(id)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch article by ID from repository")
		return nil, errors.New("article not found")
	}

	article.AuthorID = input.AuthorID
	article.Title = input.Title
	article.Text = input.Text
	article.Published = input.Published

	if err := s.repo.Update(article); err != nil {
		s.Logger.WithError(err).Error("Failed to update article in repository")
		return nil, err
	}

	return article, nil
}

// DeleteArticle удаляет статью по ID.
func (s *ArticleService) DeleteArticle(id uint) error {
	s.Logger.WithField("article_id", id).Info("Deleting article in service")

	if err := s.repo.Delete(id); err != nil {
		s.Logger.WithError(err).Error("Failed to delete article from repository")
		return err
	}
	return nil
}
