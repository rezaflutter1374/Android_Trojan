package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		if strings.HasPrefix(c.Request.URL.Path, "/ws") {
			c.Next()
			return
		}
		apiKey := c.GetHeader("X-API-Key")

		if apiKey != "test-api-key" {
			log.Printf("Unauthorized access attempt from %s", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log the request
		log.Printf("%s - %s %s %s", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Request.UserAgent())

		c.Next()
		log.Printf("%s - %s %s - %d", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	}
}
