package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyeinsoe26/indego-app/internal/app/models"
	m "github.com/nyeinsoe26/indego-app/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test FetchIndegoDataAndStore - Success
func TestFetchIndegoDataAndStore_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock the services
	mockIndegoService := new(m.MockIndegoService)
	mockWeatherService := new(m.MockWeatherService)

	// Mock data
	indegoData := models.IndegoData{LastUpdated: time.Now()}
	weatherData := models.WeatherData{}

	// Set up mock expectations
	mockIndegoService.On("GetIndegoData").Return(indegoData, nil)
	mockWeatherService.On("GetWeatherData", 39.9526, -75.1652).Return(weatherData, nil)
	mockIndegoService.On("StoreIndegoData", indegoData).Return(1, nil)
	mockWeatherService.On("StoreWeatherData", weatherData).Return(1, nil)
	mockIndegoService.On("StoreSnapshotLink", 1, 1, indegoData.LastUpdated).Return(nil)

	// Set up the handler
	handler := NewHandler(mockIndegoService, mockWeatherService)

	// Perform the test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.FetchIndegoDataAndStore(c)

	// Assert that the status code is OK
	assert.Equal(t, http.StatusCreated, w.Code)
	mockIndegoService.AssertExpectations(t)
	mockWeatherService.AssertExpectations(t)
}

// Test FetchIndegoDataAndStore - Error from Indego
func TestFetchIndegoDataAndStore_IndegoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockIndegoService := new(m.MockIndegoService)
	mockWeatherService := new(m.MockWeatherService)

	mockIndegoService.On("GetIndegoData").Return(models.IndegoData{}, errors.New("failed to fetch Indego data"))

	handler := NewHandler(mockIndegoService, mockWeatherService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.FetchIndegoDataAndStore(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockIndegoService.AssertExpectations(t)
}

// Test GetStationSnapshot - Success
func TestGetStationSnapshot_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockIndegoService := new(m.MockIndegoService)
	mockWeatherService := new(m.MockWeatherService)

	indegoData := models.IndegoData{LastUpdated: time.Now()}
	weatherData := models.WeatherData{}
	mockIndegoService.On("GetSnapshot", mock.Anything).Return(indegoData, weatherData, nil)

	handler := NewHandler(mockIndegoService, mockWeatherService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/stations?at=2019-09-01T10:00:00Z", nil)
	handler.GetStationSnapshot(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockIndegoService.AssertExpectations(t)
}

// Test GetStationSnapshot - Empty 'at' Query Parameter
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
}

// Test GetSpecificStationSnapshot - Station Not Found
func TestGetSpecificStationSnapshot_StationNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockIndegoService := new(m.MockIndegoService)
	mockWeatherService := new(m.MockWeatherService)

	indegoData := models.IndegoData{}
	weatherData := models.WeatherData{}
	mockIndegoService.On("GetSnapshot", mock.Anything).Return(indegoData, weatherData, nil)

	handler := NewHandler(mockIndegoService, mockWeatherService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{{Key: "kioskId", Value: "3005"}}
	c.Request = httptest.NewRequest("GET", "/stations/3005?at=2019-09-01T10:00:00Z", nil)
	handler.GetSpecificStationSnapshot(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockIndegoService.AssertExpectations(t)
}

// Test GetSpecificStationSnapshot - Missing KioskId
func TestGetSpecificStationSnapshot_MissingKioskId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockIndegoService := new(m.MockIndegoService)
	mockWeatherService := new(m.MockWeatherService)

	handler := NewHandler(mockIndegoService, mockWeatherService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Request without kioskId
	c.Request = httptest.NewRequest("GET", "/stations/", nil)
	handler.GetSpecificStationSnapshot(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test GetSpecificStationSnapshot - Invalid KioskId
func TestGetSpecificStationSnapshot_InvalidKioskId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockIndegoService := new(m.MockIndegoService)
	mockWeatherService := new(m.MockWeatherService)

	handler := NewHandler(mockIndegoService, mockWeatherService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Request with invalid kioskId (non-numeric)
	c.Request = httptest.NewRequest("GET", "/stations/abc?at=2019-09-01T10:00:00Z", nil)
	handler.GetSpecificStationSnapshot(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test FetchIndegoDataAndStore - Error from WeatherService
func TestFetchIndegoDataAndStore_WeatherError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockIndegoService := new(m.MockIndegoService)
	mockWeatherService := new(m.MockWeatherService)

	indegoData := models.IndegoData{LastUpdated: time.Now()}

	mockIndegoService.On("GetIndegoData").Return(indegoData, nil)
	mockWeatherService.On("GetWeatherData", 39.9526, -75.1652).Return(models.WeatherData{}, errors.New("weather fetch error"))

	handler := NewHandler(mockIndegoService, mockWeatherService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.FetchIndegoDataAndStore(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockIndegoService.AssertExpectations(t)
	mockWeatherService.AssertExpectations(t)
}

// Test GetStationSnapshot - Invalid Date Format
func TestGetStationSnapshot_InvalidDateFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockIndegoService := new(m.MockIndegoService)
	mockWeatherService := new(m.MockWeatherService)

	handler := NewHandler(mockIndegoService, mockWeatherService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Invalid date format
	c.Request = httptest.NewRequest("GET", "/stations?at=invalid-date-format", nil)
	handler.GetStationSnapshot(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test GetSpecificStationSnapshot - Success
func TestGetSpecificStationSnapshot_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockIndegoService := new(m.MockIndegoService)
	mockWeatherService := new(m.MockWeatherService)

	indegoData := models.IndegoData{
		Features: []models.StationFeature{
			{
				Properties: models.StationProperties{ID: 3005},
			},
		},
		LastUpdated: time.Now(),
	}
	weatherData := models.WeatherData{}

	mockIndegoService.On("GetSnapshot", mock.Anything).Return(indegoData, weatherData, nil)

	handler := NewHandler(mockIndegoService, mockWeatherService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{{Key: "kioskId", Value: "3005"}}
	c.Request = httptest.NewRequest("GET", "/stations/3005?at=2019-09-01T10:00:00Z", nil)
	handler.GetSpecificStationSnapshot(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockIndegoService.AssertExpectations(t)
}
