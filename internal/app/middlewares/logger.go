package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs incoming requests and response times
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Process request
		c.Next()

		// Calculate how long the request took
		latency := time.Since(t)
		fmt.Printf("Request processed in %s\n", latency)
	}
}
