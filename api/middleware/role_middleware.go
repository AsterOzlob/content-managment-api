package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrForbidden = errors.New("forbidden")

// RoleMiddleware проверяет, имеет ли пользователь одну из разрешённых ролей.
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRolesI, exists := c.Get("userRoles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: no roles provided"})
			return
		}

		userRoles, ok := userRolesI.([]string)
		if !ok || len(userRoles) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid role format"})
			return
		}

		for _, allowed := range allowedRoles {
			for _, role := range userRoles {
				if role == allowed {
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
	}
}
