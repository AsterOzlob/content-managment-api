package seeds

import (
	"fmt"

	"github.com/AsterOzlob/content_managment_api/internal/database/models"
	"github.com/AsterOzlob/content_managment_api/pkg/utils"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	// Заполнение ролей
	roles := []models.Role{
		{Name: "admin", Description: "Полный доступ ко всему"},
		{Name: "author", Description: "Может писать и управлять своими статьями"},
		{Name: "moderator", Description: "Может удалять и редактировать комментарии и статьи"},
		{Name: "user", Description: "Чтение статей и оставление комментариев"},
	}

	for _, role := range roles {
		var existing models.Role
		if err := db.Where("name = ?", role.Name).First(&existing).Error; err != nil {
			if err := db.Create(&role).Error; err != nil {
				return err
			}
		}
	}

	// Создание пользователей
	users := []models.User{
		{
			Username: "admin",
			Email:    "admin@example.com",
			RoleID:   1,
		},
		{
			Username: "john_doe",
			Email:    "john@example.com",
			RoleID:   2,
		},
		{
			Username: "jane_moderator",
			Email:    "jane@example.com",
			RoleID:   3,
		},
		{
			Username: "guest_user",
			Email:    "guest@example.com",
			RoleID:   4,
		},
	}

	// Хэширование паролей
	for i, user := range users {
		hashedPassword, err := utils.HashPassword("password")
		if err != nil {
			return fmt.Errorf("failed to hash password for user %s: %w", user.Username, err)
		}
		users[i].PasswordHash = hashedPassword
	}

	// Сохранение пользователей
	for _, user := range users {
		var existing models.User
		if err := db.Where("username = ?", user.Username).First(&existing).Error; err != nil {
			if err := db.Create(&user).Error; err != nil {
				return err
			}
		}
	}

	// Статьи
	articles := []models.Article{
		{
			Title:     "Как начать программировать",
			Text:      "Программирование — это искусство создания решений через код...",
			AuthorID:  2,
			Published: true,
		},
		{
			Title:     "Введение в Golang",
			Text:      "Go — это язык программирования, созданный Google...",
			AuthorID:  2,
			Published: true,
		},
		{
			Title:     "Работа с базами данных",
			Text:      "Базы данных — основа любого приложения...",
			AuthorID:  3,
			Published: false,
		},
	}
	for _, article := range articles {
		var existing models.Article
		if err := db.Where("title = ?", article.Title).First(&existing).Error; err != nil {
			if err := db.Create(&article).Error; err != nil {
				return err
			}
		}
	}

	// Комментарии
	comments := []models.Comment{
		{
			Text:      "Отличная статья!",
			ArticleID: 1,
			AuthorID:  4,
		},
		{
			Text:      "Мне понравилось объяснение.",
			ArticleID: 1,
			AuthorID:  3,
		},
		{
			Text:      "А как насчёт примеров кода?",
			ArticleID: 1,
			AuthorID:  4,
			ParentID:  &[]uint{1}[0],
		},
		{
			Text:      "Хороший старт!",
			ArticleID: 2,
			AuthorID:  1,
		},
	}
	for _, comment := range comments {
		var existing models.Comment
		if err := db.Where("text = ?", comment.Text).First(&existing).Error; err != nil {
			if err := db.Create(&comment).Error; err != nil {
				return err
			}
		}
	}

	// Медиафайлы
	media := []models.Media{
		{
			FilePath:  "/uploads/go-logo.png",
			FileType:  "image/png",
			FileSize:  10240,
			AuthorID:  2,
			ArticleID: &[]uint{1}[0],
		},
		{
			FilePath:  "/uploads/db-diagram.jpg",
			FileType:  "image/jpeg",
			FileSize:  20480,
			AuthorID:  3,
			ArticleID: &[]uint{3}[0],
		},
	}
	for _, m := range media {
		var existing models.Media
		if err := db.Where("file_path = ?", m.FilePath).First(&existing).Error; err != nil {
			if err := db.Create(&m).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
