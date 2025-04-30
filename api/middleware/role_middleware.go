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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
			return
		}

		roleStr := role.(string)
		allowed := false
		for _, r := range allowedRoles {
			if roleStr == r {
				allowed = true
				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": ErrForbidden.Error()})
			return
		}

		c.Next()
	}
}
