package api

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/nyeinsoe26/indego-app/internal/app/api/callback"
	"github.com/nyeinsoe26/indego-app/internal/app/api/login"
	"github.com/nyeinsoe26/indego-app/internal/app/api/logout"
	"github.com/nyeinsoe26/indego-app/internal/app/api/user"
	m "github.com/nyeinsoe26/indego-app/internal/app/middlewares"
	"github.com/nyeinsoe26/indego-app/internal/app/middlewares/authenticator"
)

func RegisterRoutes(router *gin.Engine, handler *Handler, auth *authenticator.Authenticator) {
	// Create a cookie-based session store
	store := cookie.NewStore([]byte("super-secret-key"))
	router.Use(sessions.Sessions("auth-session", store))

	// Logger middleware
	router.Use(m.LoggerMiddleware())

	// Public routes
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth routes
	router.GET("/login", login.Handler(auth))
	router.GET("/callback", callback.Handler(auth))
	router.GET("/logout", logout.Handler)

	// User route - protected
	router.GET("/user", m.IsAuthenticated, user.Handler)

	// API routes - protected
	v1 := router.Group("/api/v1")
	v1.Use(m.IsAuthenticated)
	{
		v1.POST("/indego-data-fetch-and-store-it-db", handler.FetchIndegoDataAndStore)
		v1.GET("/stations", handler.GetStationSnapshot)
		v1.GET("/stations/:kioskId", handler.GetSpecificStationSnapshot)
	}
}
