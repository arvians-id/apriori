package middleware

import (
	"github.com/gin-gonic/gin"
	"os"
)

func SetupCorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-KEY")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, GET, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func SetupXApiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey != os.Getenv("X_API_KEY") {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.AbortWithStatus(401)
			return
		}

		c.Next()
	}
}
