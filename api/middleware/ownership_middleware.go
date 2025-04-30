package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OwnershipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем ID пользователя из контекста
		userID, userIDExists := c.Get("userID")
		if !userIDExists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}

		// Извлекаем ID владельца ресурса из URL
		resourceOwnerID := c.Param("id")
		if resourceOwnerID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Resource owner ID not provided"})
			return
		}

		// Преобразуем userID к строке для сравнения с resourceOwnerID
		userIDStr := fmt.Sprintf("%d", userID.(uint))
		if userIDStr != resourceOwnerID {
			// Если пользователь не является владельцем, проверяем роль
			role, exists := c.Get("role")
			if !exists || (role != "admin" && role != "moderator") {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this resource"})
				return
			}
		}

		c.Next()
	}
}
