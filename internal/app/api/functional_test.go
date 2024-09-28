package api_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyeinsoe26/indego-app/internal/app/api"
	"github.com/nyeinsoe26/indego-app/internal/app/models"
	"github.com/nyeinsoe26/indego-app/internal/app/services"
	"github.com/nyeinsoe26/indego-app/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var ts *httptest.Server
var mockDB *mocks.MockDatabase
var mockIndegoClient *mocks.MockIndegoClient
var mockWeatherClient *mocks.MockWeatherClient

// setupTestRouter sets up the Gin router for testing with mock services and mock DB.
func setupTestRouter() *gin.Engine {
	// Initialize Gin router
	router := gin.Default()

	// Create services with mock clients and mock DB
	indegoService := services.NewIndegoService(mockIndegoClient, mockDB)
	weatherService := services.NewWeatherService(mockWeatherClient, mockDB)

	// Initialize the handler with mock services and mock DB
	handler := api.NewHandler(indegoService, weatherService)

	// Register routes without the authentication middleware for testing
	v1 := router.Group("/api/v1")
	{
		v1.POST("/indego-data-fetch-and-store-it-db", handler.FetchIndegoDataAndStore)
		v1.GET("/stations", handler.GetStationSnapshot)
		v1.GET("/stations/:kioskId", handler.GetSpecificStationSnapshot)
	}

	return router
}

// TestMain is used to set up global resources before running the tests
func TestMain(m *testing.M) {
	// Set up the mock database and clients
	mockDB = new(mocks.MockDatabase)
	mockIndegoClient = new(mocks.MockIndegoClient)
	mockWeatherClient = new(mocks.MockWeatherClient)

	// Set up the real router once
	router := setupTestRouter()

	// Start the test HTTP server once
	ts = httptest.NewServer(router)
	defer ts.Close()

	// Run all tests
	code := m.Run()

	// Clean up resources after tests are done
	os.Exit(code)
}

// Helper function to get the full URL for the test server
func getTestURL(path string) string {
	return ts.URL + path
}

// Reset mocks before each test case
func resetMocks() {
	mockDB.ExpectedCalls = nil
	mockIndegoClient.ExpectedCalls = nil
	mockWeatherClient.ExpectedCalls = nil

	mockDB.Calls = nil
	mockIndegoClient.Calls = nil
	mockWeatherClient.Calls = nil
}

// TestFunctional_FetchIndegoDataAndStore_Success tests the full flow of fetching and storing Indego
// and weather data in the database, returning a success response when everything works.
func TestFunctional_FetchIndegoDataAndStore_Success(t *testing.T) {
	resetMocks()

	// Mock data
	indegoData := models.IndegoData{LastUpdated: time.Now()}
	weatherData := models.WeatherData{}

	// Set up mock client and DB expectations
	mockIndegoClient.On("FetchIndegoData").Return(indegoData, nil)
	mockWeatherClient.On("FetchWeatherData", 39.9526, -75.1652).Return(weatherData, nil)
	mockDB.On("StoreIndegoData", indegoData).Return(1, nil)
	mockDB.On("StoreWeatherData", weatherData).Return(1, nil)
	mockDB.On("StoreSnapshotLink", 1, 1, indegoData.LastUpdated).Return(nil)

	// Make the POST request to the actual endpoint
	resp, err := http.Post(getTestURL("/api/v1/indego-data-fetch-and-store-it-db"), "application/json", nil)
	assert.NoError(t, err)

	// Assert that the status code is 201 Created
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	mockDB.AssertExpectations(t)
	mockIndegoClient.AssertExpectations(t)
	mockWeatherClient.AssertExpectations(t)
}

// TestFunctional_FetchIndegoDataAndStore_IndegoError tests the case where fetching Indego data fails,
// ensuring the handler doesn't proceed to weather fetching or database storage, and returns a 500 error.
func TestFunctional_FetchIndegoDataAndStore_IndegoError(t *testing.T) {
	resetMocks()

	// Mock IndegoClient to return an error
	mockIndegoClient.On("FetchIndegoData").Return(models.IndegoData{}, assert.AnError)

	// Ensure no calls are made to WeatherClient or DB when IndegoClient fails
	mockWeatherClient.AssertNotCalled(t, "FetchWeatherData", mock.Anything, mock.Anything)
	mockDB.AssertNotCalled(t, "StoreIndegoData", mock.Anything)
	mockDB.AssertNotCalled(t, "StoreWeatherData", mock.Anything)
	mockDB.AssertNotCalled(t, "StoreSnapshotLink", mock.Anything, mock.Anything, mock.Anything)

	// Make the POST request to the actual endpoint
	resp, err := http.Post(getTestURL("/api/v1/indego-data-fetch-and-store-it-db"), "application/json", nil)
	assert.NoError(t, err)

	// Assert that the status code is 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	// Verify that the mocks were called as expected
	mockIndegoClient.AssertExpectations(t)
}

// TestFuncional_GetStationSnapshot_Success tests fetching station and weather snapshot data
// for a valid 'at' query parameter, expecting a 200 OK response and correct data from the DB.
func TestFuncional_GetStationSnapshot_Success(t *testing.T) {
	resetMocks()

	// Mock data
	indegoData := models.IndegoData{LastUpdated: time.Now()}
	weatherData := models.WeatherData{}

	// Set up mock DB expectations
	mockDB.On("FetchSnapshot", mock.Anything).Return(indegoData, weatherData, nil)

	// Make the GET request with a valid query parameter
	resp, err := http.Get(getTestURL("/api/v1/stations?at=2019-09-01T10:00:00Z"))
	assert.NoError(t, err)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockDB.AssertExpectations(t)
}

// TestFunctional_GetStationSnapshot_InvalidDateFormat tests that if the 'at' query parameter
// is invalid, the handler returns a 400 Bad Request response.
func TestFunctional_GetStationSnapshot_InvalidDateFormat(t *testing.T) {
	resetMocks()

	// Make the GET request with an invalid date format
	resp, err := http.Get(getTestURL("/api/v1/stations?at=invalid-date-format"))
	assert.NoError(t, err)

	// Assert that the status code is 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestFunctional_GetSpecificStationSnapshot_Success tests fetching a specific station snapshot
// by its kioskId and time, ensuring a 200 OK response if the station is found.
func TestFunctional_GetSpecificStationSnapshot_Success(t *testing.T) {
	resetMocks()

	// Mock data for a station with kioskId 3005
	indegoData := models.IndegoData{
		Features: []models.StationFeature{
			{
				Properties: models.StationProperties{ID: 3005},
			},
		},
		LastUpdated: time.Now(),
	}
	weatherData := models.WeatherData{}

	// Set up mock DB expectations
	mockDB.On("FetchSnapshot", mock.Anything).Return(indegoData, weatherData, nil)

	// Make the GET request with the kioskId and a query parameter
	resp, err := http.Get(getTestURL("/api/v1/stations/3005?at=2019-09-01T10:00:00Z"))
	assert.NoError(t, err)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockDB.AssertExpectations(t)
}

// TestFunctional_GetSpecificStationSnapshot_InvalidKioskId tests the scenario where an invalid (non-numeric)
// kioskId is passed, and the handler should return a 400 Bad Request response.
func TestFunctional_GetSpecificStationSnapshot_InvalidKioskId(t *testing.T) {
	resetMocks()

	// Make the GET request with an invalid kioskId (non-numeric)
	resp, err := http.Get(getTestURL("/api/v1/stations/abc?at=2019-09-01T10:00:00Z"))
	assert.NoError(t, err)

	// Assert that the status code is 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestFunctional_GetSpecificStationSnapshot_StationNotFound tests the case where a specific station
// snapshot is requested but the station is not found in the data, expecting a 404 Not Found response.
func TestFunctional_GetSpecificStationSnapshot_StationNotFound(t *testing.T) {
	resetMocks()

	// Mock data
	indegoData := models.IndegoData{}
	weatherData := models.WeatherData{}

	// Set up mock DB expectations
	mockDB.On("FetchSnapshot", mock.Anything).Return(indegoData, weatherData, nil)

	// Make the GET request with the kioskId and a query parameter
	resp, err := http.Get(getTestURL("/api/v1/stations/3005?at=2019-09-01T10:00:00Z"))
	assert.NoError(t, err)

	// Assert that the status code is 404 Not Found
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	mockDB.AssertExpectations(t)
}
