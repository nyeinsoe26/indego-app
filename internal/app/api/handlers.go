package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyeinsoe26/indego-app/internal/app/services"
)

// Handler struct with services
type Handler struct {
	IndegoService  services.IndegoService
	WeatherService services.WeatherService
}

// NewHandler creates a new Handler with Indego and Weather services
func NewHandler(indegoService services.IndegoService, weatherService services.WeatherService) *Handler {
	return &Handler{
		IndegoService:  indegoService,
		WeatherService: weatherService,
	}
}

// Core logic to fetch and store Indego and weather data, independent of gin.Context
func (h *Handler) FetchAndStoreIndegoWeatherData() error {
	// Fetch Indego data using the Indego service
	indegoData, err := h.IndegoService.GetIndegoData()
	if err != nil {
		return fmt.Errorf("failed to fetch Indego data: %w", err)
	}

	// Fetch Weather data using the Weather service
	weatherData, err := h.WeatherService.GetWeatherData(39.9526, -75.1652)
	if err != nil {
		return fmt.Errorf("failed to fetch weather data: %w", err)
	}

	// Store Indego data in the database
	err = h.IndegoService.StoreIndegoData(indegoData)
	if err != nil {
		return fmt.Errorf("failed to store Indego data: %w", err)
	}

	// Store Weather data in the database
	err = h.WeatherService.StoreWeatherData(weatherData)
	if err != nil {
		return fmt.Errorf("failed to store weather data: %w", err)
	}

	return nil
}

// FetchIndegoDataAndStore stores the latest Indego and Weather data in the database
func (h *Handler) FetchIndegoDataAndStore(c *gin.Context) {
	// Call the core logic
	err := h.FetchAndStoreIndegoWeatherData()
	if err != nil {
		// Respond with an error message if something goes wrong
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a success message if everything goes well
	c.JSON(http.StatusOK, gin.H{"message": "Data stored successfully"})
}

// GetStationSnapshot retrieves a snapshot of all stations and weather data at a specific time
func (h *Handler) GetStationSnapshot(c *gin.Context) {
	// Parse the 'at' parameter from the query string
	timeStr := c.Query("at")
	at, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format"})
		return
	}

	// Fetch the snapshot from the service
	indegoData, weatherData, err := h.IndegoService.GetSnapshot(at)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch snapshot"})
		return
	}

	// Respond with the snapshot data
	c.JSON(http.StatusOK, gin.H{
		"at":       at.Format(time.RFC3339),
		"stations": indegoData,
		"weather":  weatherData,
	})
}

// GetSpecificStationSnapshot retrieves a snapshot of a specific station at a specific time
func (h *Handler) GetSpecificStationSnapshot(c *gin.Context) {
	// Extract the kioskId from the URL parameters
	kioskIDStr := c.Param("kioskId")

	// Convert kioskID from string to int
	kioskID, err := strconv.Atoi(kioskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid kioskId format"})
		return
	}

	// Parse the 'at' parameter from the query string
	timeStr := c.Query("at")
	at, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format"})
		return
	}

	// Fetch the snapshot from the service
	indegoData, weatherData, err := h.IndegoService.GetSnapshot(at)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch snapshot"})
		return
	}

	// Find the specific station by kioskId
	var stationData interface{}
	for _, station := range indegoData.Features {
		if station.Properties.ID == kioskID {
			stationData = station
			break
		}
	}

	// If the station is not found, return 404
	if stationData == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Station not found"})
		return
	}

	// Respond with the station snapshot data
	c.JSON(http.StatusOK, gin.H{
		"at":      at.Format(time.RFC3339),
		"station": stationData,
		"weather": weatherData,
	})
}
