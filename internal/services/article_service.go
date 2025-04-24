package services

import (
	"errors"

	"github.com/AsterOzlob/content_managment_api/internal/dto"
	"github.com/AsterOzlob/content_managment_api/internal/models"
	"github.com/AsterOzlob/content_managment_api/internal/repositories"
	logging "github.com/AsterOzlob/content_managment_api/logger"
	"github.com/sirupsen/logrus"
)

// ArticleService предоставляет методы для управления статьями.
type ArticleService struct {
	repo   *repositories.ArticleRepository // repo - репозиторий для взаимодействия с базой данных.
	Logger *logging.Logger                 // Logger - экземпляр логгера для ArticleService.
}

// NewArticleService создает новый экземпляр ArticleService.
func NewArticleService(repo *repositories.ArticleRepository, logger *logging.Logger) *ArticleService {
	return &ArticleService{repo: repo, Logger: logger}
}

// CreateArticle создает новую статью.
func (s *ArticleService) CreateArticle(input dto.ArticleInput) (*models.Article, error) {
	s.Logger.Log(logrus.InfoLevel, "Creating article in service", map[string]interface{}{
		"author_id": input.AuthorID,
		"title":     input.Title,
	})

	// Проверяем, существует ли пользователь с указанным AuthorID
	var user models.User
	if err := s.repo.DB.First(&user, input.AuthorID).Error; err != nil {
		s.Logger.Log(logrus.ErrorLevel, "User not found", map[string]interface{}{
			"error":     err.Error(),
			"author_id": input.AuthorID,
		})
		return nil, errors.New("user not found")
	}

	article := &models.Article{
		AuthorID:  input.AuthorID,
		Title:     input.Title,
		Text:      input.Text,
		Published: input.Published,
	}

	if err := s.repo.Create(article); err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to create article in repository", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	if len(input.MediaIDs) > 0 {
		if err := s.repo.AttachMedia(article.ID, input.MediaIDs); err != nil {
			s.Logger.Log(logrus.ErrorLevel, "Failed to attach media to article", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}
	}

	return article, nil
}

// GetAllArticles возвращает список всех статей.
func (s *ArticleService) GetAllArticles() ([]*models.Article, error) {
	s.Logger.Log(logrus.InfoLevel, "Fetching all articles in service", nil)

	articles, err := s.repo.GetAll()
	if err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to fetch all articles from repository", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	return articles, nil
}

// GetArticleByID возвращает статью по ID.
func (s *ArticleService) GetArticleByID(id uint) (*models.Article, error) {
	s.Logger.Log(logrus.InfoLevel, "Fetching article by ID in service", map[string]interface{}{
		"article_id": id,
	})

	article, err := s.repo.GetByID(id)
	if err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to fetch article by ID from repository", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("article not found")
	}
	return article, nil
}

// UpdateArticle обновляет существующую статью.
func (s *ArticleService) UpdateArticle(id uint, input dto.ArticleInput) (*models.Article, error) {
	s.Logger.Log(logrus.InfoLevel, "Updating article in service", map[string]interface{}{
		"article_id": id,
		"title":      input.Title,
	})

	article, err := s.repo.GetByID(id)
	if err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to fetch article by ID from repository", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("article not found")
	}

	article.AuthorID = input.AuthorID
	article.Title = input.Title
	article.Text = input.Text
	article.Published = input.Published

	if err := s.repo.Update(article); err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to update article in repository", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	// Обновляем медиафайлы, если они указаны.
	if len(input.MediaIDs) > 0 {
		if err := s.repo.AttachMedia(article.ID, input.MediaIDs); err != nil {
			s.Logger.Log(logrus.ErrorLevel, "Failed to attach media to article", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}
	}

	return article, nil
}

// DeleteArticle удаляет статью по ID.
func (s *ArticleService) DeleteArticle(id uint) error {
	s.Logger.Log(logrus.InfoLevel, "Deleting article in service", map[string]interface{}{
		"article_id": id,
	})

	if err := s.repo.Delete(id); err != nil {
		s.Logger.Log(logrus.ErrorLevel, "Failed to delete article from repository", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	return nil
}
