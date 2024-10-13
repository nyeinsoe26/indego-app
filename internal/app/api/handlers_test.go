package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyeinsoe26/indego-app/internal/app/dtos"
	"github.com/nyeinsoe26/indego-app/internal/app/models"
	m "github.com/nyeinsoe26/indego-app/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestFetchIndegoDataAndStore_Success tests the happy path where both Indego and Weather data
// are successfully fetched and stored, and the correct response is returned.
func TestFetchIndegoDataAndStore_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock the services
	mockIndegoService := new(m.MockIndegoService)
	mockWeatherService := new(m.MockWeatherService)

	// Mock data
	indegoData := models.IndegoData{LastUpdated: time.Now()}
	weatherData := models.WeatherData{}
	mockUUID := uuid.New()

	// Set up mock expectations
	mockIndegoService.On("GetIndegoData").Return(indegoData, nil)
	mockWeatherService.On("GetWeatherData", 39.9526, -75.1652).Return(weatherData, nil)
	mockIndegoService.On("StoreSnapshot", indegoData, weatherData, indegoData.LastUpdated).Return(mockUUID, nil)

	// Set up the handler
	handler := NewHandler(mockIndegoService, mockWeatherService)

	// Perform the test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.FetchIndegoDataAndStore(c)

	// Assert that the status code is Created (201)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Assert that the response is the expected DTO
	var response dtos.FetchIndegoWeatherResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Data stored successfully", response.Message)

	mockIndegoService.AssertExpectations(t)
	mockWeatherService.AssertExpectations(t)
}

// TestFetchIndegoDataAndStore_IndegoError tests the scenario where the Indego data
// fetching fails, and the handler should return an error after retry attempts.
func TestFetchIndegoDataAndStore_IndegoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockIndegoService := new(m.MockIndegoService)
	mockWeatherService := new(m.MockWeatherService)

	// The mock should simulate the Indego service returning an error
	mockIndegoService.On("GetIndegoData").Return(models.IndegoData{}, errors.New("failed to fetch Indego data"))

	handler := NewHandler(mockIndegoService, mockWeatherService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Trigger the handler
	handler.FetchIndegoDataAndStore(c)

	// Assert that the status code is Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// The expected error message now includes the retry attempts and the detailed error
	expectedErrorMessage := "failed to fetch Indego data after 3 attempts: failed to fetch Indego data"

	// Assert that the response contains the correct error message
	var errorResponse dtos.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedErrorMessage, errorResponse.Error)

	mockIndegoService.AssertExpectations(t)
}

// TestGetStationSnapshot_Success tests the scenario where a station snapshot is
// successfully retrieved and the response contains the correct data.
func TestGetStationSnapshot_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock the services
	mockIndegoService := new(m.MockIndegoService)
	mockWeatherService := new(m.MockWeatherService)

	// Mock data
	indegoData := models.IndegoData{
		LastUpdated: time.Date(2024, time.September, 28, 18, 47, 50, 0, time.UTC),
		Features:    []models.StationFeature{},
	}
	weatherData := models.WeatherData{}
	snapshotTime := time.Date(2024, time.September, 28, 18, 47, 50, 0, time.UTC)

	// Set up mock expectations
	mockIndegoService.On("GetSnapshot", mock.Anything).Return(indegoData, weatherData, snapshotTime, nil)

	// Set up the handler
	handler := NewHandler(mockIndegoService, mockWeatherService)

	// Perform the test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/stations?at=2024-09-28T18:47:50Z", nil)
	handler.GetStationSnapshot(c)

	// Assert that the status code is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Extract the response data
	var response dtos.StationSnapshotResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Compare the expected and actual LastUpdated times
	expectedTime := snapshotTime
	actualTime, err := time.Parse(time.RFC3339, response.At)
	assert.NoError(t, err)

	assert.Equal(t, expectedTime, actualTime)
	assert.Equal(t, indegoData.Features, response.Stations.Features)

	mockIndegoService.AssertExpectations(t)
}

// TestGetStationSnapshot_EmptyAtParam tests the scenario where the 'at' query parameter
// is missing or empty, and the handler should return a 400 Bad Request.
func TestGetStationSnapshot_EmptyAtParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockIndegoService := new(m.MockIndegoService)
	mockWeatherService := new(m.MockWeatherService)

	handler := NewHandler(mockIndegoService, mockWeatherService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Request with empty 'at' parameter
	c.Request = httptest.NewRequest("GET", "/stations?at=", nil)
	handler.GetStationSnapshot(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Assert that the error response is returned
	var errorResponse dtos.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid time format", errorResponse.Error)
}

// TestGetSpecificStationSnapshot_StationNotFound tests the case where the requested station
// is not found in the snapshot, and the handler should return a 404 Not Found.
func TestGetSpecificStationSnapshot_StationNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockIndegoService := new(m.MockIndegoService)
	mockWeatherService := new(m.MockWeatherService)

	indegoData := models.IndegoData{}
	weatherData := models.WeatherData{}
	snapshotTime := time.Now()

	mockIndegoService.On("GetSnapshot", mock.Anything).Return(indegoData, weatherData, snapshotTime, nil)

	handler := NewHandler(mockIndegoService, mockWeatherService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{{Key: "kioskId", Value: "3005"}}
	c.Request = httptest.NewRequest("GET", "/stations/3005?at=2019-09-01T10:00:00Z", nil)
	handler.GetSpecificStationSnapshot(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

	// Assert that the error response is returned
	var errorResponse dtos.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Station not found", errorResponse.Error)

	mockIndegoService.AssertExpectations(t)
}

// TestFetchIndegoDataAndStore_WeatherError tests the scenario where the Weather data fetching
// fails, and the handler should return an error.
func TestFetchIndegoDataAndStore_WeatherError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock the services
	mockIndegoService := new(m.MockIndegoService)
	mockWeatherService := new(m.MockWeatherService)

	indegoData := models.IndegoData{LastUpdated: time.Now()}

	// Set up mock expectations
	mockIndegoService.On("GetIndegoData").Return(indegoData, nil)
	mockWeatherService.On("GetWeatherData", 39.9526, -75.1652).Return(models.WeatherData{}, errors.New("weather fetch error"))

	// Set up the handler
	handler := NewHandler(mockIndegoService, mockWeatherService)

	// Perform the test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.FetchIndegoDataAndStore(c)

	// Assert that the status code is 500
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert that the error message contains the expected substring
	assert.Contains(t, w.Body.String(), "weather fetch error")

	mockIndegoService.AssertExpectations(t)
	mockWeatherService.AssertExpectations(t)
}
