package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckOwnershipMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID") // ID пользователя из JWT
		resourceOwnerID := c.Param("id")

		if userID != resourceOwnerID {
			for _, role := range allowedRoles {
				if role == "admin" && c.GetString("role") == "admin" {
					return
				}
			}
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this resource"})
			return
		}
		c.Next()
	}
}
