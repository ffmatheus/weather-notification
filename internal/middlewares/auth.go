package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiToken := os.Getenv("API_TOKEN")
		if apiToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API token não configurado"})
			c.Abort()
			return
		}

		requestToken := c.GetHeader("Authorization")
		if requestToken == "" || requestToken != "Bearer "+apiToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou ausente"})
			c.Abort()
			return
		}

		c.Next()
	}
}
