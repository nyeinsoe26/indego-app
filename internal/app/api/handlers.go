package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyeinsoe26/indego-app/internal/app/dtos"
	"github.com/nyeinsoe26/indego-app/internal/app/models"
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

// FetchAndStoreIndegoWeatherData handles the core logic for fetching and storing Indego and weather data
func (h *Handler) FetchAndStoreIndegoWeatherData() error {
	const maxRetries = 3
	var err error
	var indegoData models.IndegoData

	// Retry loop for fetching Indego data
	for i := 0; i < maxRetries; i++ {
		indegoData, err = h.IndegoService.GetIndegoData()
		if err == nil {
			break
		}
		log.Printf("Failed to fetch Indego data, attempt %d/%d: %v\n", i+1, maxRetries, err)
		time.Sleep(2 * time.Second) // Delay before retrying
	}
	if err != nil {
		// Return early if Indego data fetch fails
		return fmt.Errorf("failed to fetch Indego data after %d attempts: %w", maxRetries, err)
	}

	// Fetch Weather data using the Weather service
	weatherData, err := h.WeatherService.GetWeatherData(39.9526, -75.1652)
	if err != nil {
		return fmt.Errorf("failed to fetch weather data: %w", err)
	}

	_, err = h.IndegoService.StoreSnapshot(indegoData, weatherData, indegoData.LastUpdated)
	if err != nil {
		return fmt.Errorf("failed to store snapshot link: %w", err)
	}
	return nil
}

// FetchIndegoDataAndStore godoc
// @Summary Store the latest Indego and Weather data
// @Description Fetch the latest data from Indego and Weather services, store them in the database, and link them.
// @Tags Indego
// @Accept  json
// @Produce  json
// @Success 201 {object} dtos.FetchIndegoWeatherResponse "Data stored successfully"
// @Failure 500 {object} dtos.ErrorResponse "Internal Server Error"
// @Router /api/v1/indego-data-fetch-and-store-it-db [post]
func (h *Handler) FetchIndegoDataAndStore(c *gin.Context) {
	// Call the core logic
	err := h.FetchAndStoreIndegoWeatherData()
	if err != nil {
		// Respond with an error message if something goes wrong
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Respond with a success message if everything goes well
	c.JSON(http.StatusCreated, dtos.FetchIndegoWeatherResponse{
		Message: "Data stored successfully",
	})
}

// GetStationSnapshot godoc
// @Summary Retrieve a snapshot of all stations at a specific time
// @Description Get a snapshot of all stations and weather data at a specified time using the 'at' query parameter.
// @Tags Indego
// @Accept  json
// @Produce  json
// @Param  at  query  string  true  "Timestamp in RFC3339 format"
// @Success 200 {object} dtos.StationSnapshotResponse "Snapshot data"
// @Failure 400 {object} dtos.ErrorResponse "Invalid time format"
// @Failure 500 {object} dtos.ErrorResponse "Internal Server Error"
// @Router /api/v1/stations [get]
func (h *Handler) GetStationSnapshot(c *gin.Context) {
	// Parse the 'at' parameter from the query string
	timeStr := c.Query("at")
	at, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Error: "Invalid time format",
		})
		return
	}

	// Fetch the snapshot from the service
	indegoData, weatherData, snapshotTime, err := h.IndegoService.GetSnapshot(at)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Error: "Failed to fetch snapshot",
		})
		return
	}

	// Respond with the snapshot data
	c.JSON(http.StatusOK, dtos.StationSnapshotResponse{
		At:       snapshotTime.Format(time.RFC3339),
		Stations: indegoData,
		Weather:  weatherData,
	})
}

// GetSpecificStationSnapshot godoc
// @Summary Retrieve a snapshot of a specific station at a specific time
// @Description Get a snapshot of a specific station's data at a given time.
// @Tags Indego
// @Accept  json
// @Produce  json
// @Param  kioskId  path  string  true  "Kiosk ID"
// @Param  at  query  string  true  "Timestamp in RFC3339 format"
// @Success 200 {object} dtos.SpecificStationSnapshotResponse "Station data"
// @Failure 400 {object} dtos.ErrorResponse "Invalid kioskId or time format"
// @Failure 404 {object} dtos.ErrorResponse "Station not found"
// @Failure 500 {object} dtos.ErrorResponse "Failed to fetch snapshot"
// @Router /api/v1/stations/{kioskId} [get]
func (h *Handler) GetSpecificStationSnapshot(c *gin.Context) {
	// Extract the kioskId from the URL parameters
	kioskIDStr := c.Param("kioskId")

	// Convert kioskID from string to int
	kioskID, err := strconv.Atoi(kioskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Error: "Invalid kioskId format",
		})
		return
	}

	// Parse the 'at' parameter from the query string
	timeStr := c.Query("at")
	at, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Error: "Invalid time format",
		})
		return
	}

	// Fetch the snapshot from the service
	indegoData, weatherData, snapshotTime, err := h.IndegoService.GetSnapshot(at)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Error: "Failed to fetch snapshot",
		})
		return
	}

	// Find the specific station by kioskId
	var stationData models.StationFeature
	found := false
	for _, station := range indegoData.Features {
		if station.Properties.ID == kioskID {
			stationData = station
			found = true
			break
		}
	}

	// If the station is not found, return 404
	if !found {
		c.JSON(http.StatusNotFound, dtos.ErrorResponse{
			Error: "Station not found",
		})
		return
	}

	// Respond with the station snapshot data
	c.JSON(http.StatusOK, dtos.SpecificStationSnapshotResponse{
		At:      snapshotTime.Format(time.RFC3339),
		Station: stationData,
		Weather: weatherData,
	})
}
