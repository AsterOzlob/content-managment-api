package middleware

import "github.com/gin-gonic/gin"

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пропускаем защитные заголовки для Swagger
		if len(c.Request.URL.Path) >= 5 && c.Request.URL.Path[:5] == "/docs" {
			c.Next()
			return
		}

		// Применяем защитные заголовки ко всем остальным маршрутам
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		c.Header("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' https://example.com ; "+
				"style-src 'self' https://fonts.googleapis.com ; "+
				"img-src 'self' data:; "+
				"font-src 'self' https://fonts.gstatic.com ; "+
				"object-src 'none'; "+
				"base-uri 'self'; "+
				"form-action 'self';")

		c.Next()
	}
}
