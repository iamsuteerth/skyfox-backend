package security

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func APIKeyAuthMiddleware() gin.HandlerFunc {
	apiKey := os.Getenv("API_GATEWAY_KEY")

	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		requestKey := c.GetHeader("X-Api-Key")

		if requestKey != apiKey {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
