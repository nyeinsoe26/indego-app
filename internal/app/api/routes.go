package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	m "github.com/nyeinsoe26/indego-app/internal/app/middlewares"
)

// RegisterRoutes sets up the API routes
func RegisterRoutes(router *gin.Engine, handler *Handler) {
	router.Use(m.LoggerMiddleware())
	router.Use(m.TokenAuthMiddleware)

	// Liveness probe route
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/indego-data-fetch-and-store-it-db", handler.FetchIndegoDataAndStore)
		v1.GET("/stations", handler.GetStationSnapshot)
		v1.GET("/stations/:kioskId", handler.GetSpecificStationSnapshot)
	}
}
