// internal/middleware/xss.go
package middleware

import (
	"html"
	"net/http"

	"github.com/gin-gonic/gin"
)

func XSSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Обрабатываем только POST, PUT, PATCH запросы
		if c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut ||
			c.Request.Method == http.MethodPatch {

			if err := c.Request.ParseForm(); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
				return
			}

			// Экранируем все поля формы
			for key, values := range c.Request.PostForm {
				for i, value := range values {
					c.Request.PostForm[key][i] = html.EscapeString(value)
				}
			}
		}

		c.Next()
	}
}
