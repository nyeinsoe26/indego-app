package middleware

import (
	"net/http"
	"strings"

	"github.com/nyeinsoe26/indego-app/config"

	"github.com/gin-gonic/gin"
)

// TokenAuthMiddleware checks for a valid Bearer token in the request headers
func TokenAuthMiddleware(c *gin.Context) {
	// Extract token from Authorization header
	token := c.GetHeader("Authorization")
	expectedToken := "Bearer " + config.AppConfig.Auth.Token

	// Validate token
	if !strings.HasPrefix(token, "Bearer ") || token != expectedToken {
		// Return 401 Unauthorized if token is invalid or missing
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort() // Stop further request processing
		return
	}

	// If the token is valid, continue with the request
	c.Next()
}
