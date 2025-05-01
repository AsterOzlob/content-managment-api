package middleware

import (
	"errors"
	"net/http"

	"github.com/AsterOzlob/content_managment_api/config"
	"github.com/AsterOzlob/content_managment_api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var ErrUnauthorized = errors.New("unauthorized")

// AuthMiddleware проверяет JWT-токен в заголовке Authorization.
func AuthMiddleware(jwtConfig *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
			return
		}

		tokenString := authHeader[7:]
		token, err := utils.ValidateAccessToken(tokenString, jwtConfig)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["user_id"].(float64))
		role := claims["role"].(string)

		// Добавляем user_id и role в контекст
		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}
