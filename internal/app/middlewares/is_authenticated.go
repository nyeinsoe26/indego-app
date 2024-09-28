package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// IsAuthenticated checks if the user is logged in

// func IsAuthenticated(ctx *gin.Context) {
// 	if sessions.Default(ctx).Get("profile") == nil {
// 		ctx.Redirect(http.StatusSeeOther, "/")
// 	} else {
// 		ctx.Next()
// 	}
// }

// CombinedAuthMiddleware handles both token and session-based authentication
func IsAuthenticated(ctx *gin.Context) {
	authHeader := ctx.Request.Header.Get("Authorization")

	// Step 1: Check if token exists in the Authorization header
	if authHeader != "" {
		// Split the token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]

			// Validate the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate the signing method
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				// Retrieve the public key
				return getPublicKey(token)
			})

			if err != nil || !token.Valid {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
				ctx.Abort()
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				ctx.Abort()
				return
			}

			// Check if the audience is valid
			if !claims.VerifyAudience(audience, true) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid audience"})
				ctx.Abort()
				return
			}

			// Token is valid, proceed with request
			ctx.Next()
			return
		}
	}

	// Step 2: If no token is found, fall back to session-based authentication
	session := sessions.Default(ctx)
	if session.Get("profile") == nil {
		// ctx.Redirect(http.StatusSeeOther, "/login")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		ctx.Abort()
		return
	}

	// Session is valid, proceed with request
	ctx.Next()
}
