package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrForbidden = errors.New("forbidden")

// RoleMiddleware проверяет, имеет ли пользователь одну из разрешенных ролей.
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		roleStr := role.(string)
		isAllowedRole := false
		for _, r := range allowedRoles {
			if roleStr == r {
				isAllowedRole = true
				break
			}
		}

		if !isAllowedRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this action"})
			return
		}

		c.Next()
	}
}
