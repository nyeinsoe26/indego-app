package services

import (
	"errors"
	"testing"

	"github.com/nyeinsoe26/indego-app/internal/app/models"
	m "github.com/nyeinsoe26/indego-app/internal/mocks"
	"github.com/stretchr/testify/assert"
)

// TestGetWeatherData_Success tests the successful retrieval of weather data
// from the mock WeatherClient, expecting no errors and correct data.
func TestGetWeatherData_Success(t *testing.T) {
	// Create mock WeatherClient
	mockWeatherClient := new(m.MockWeatherClient)
	expectedWeatherData := models.WeatherData{}
	mockWeatherClient.On("FetchWeatherData", 39.9526, -75.1652).Return(expectedWeatherData, nil)

	// Create mock Database (not used in this test)
	mockDB := new(m.MockDatabase)

	// Create the WeatherService
	service := NewWeatherService(mockWeatherClient, mockDB)

	// Call the method
	actualWeatherData, err := service.GetWeatherData(39.9526, -75.1652)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedWeatherData, actualWeatherData)
	mockWeatherClient.AssertExpectations(t)
}

// TestGetWeatherData_Error tests the scenario where fetching weather data
// from the mock WeatherClient results in an error.
func TestGetWeatherData_Error(t *testing.T) {
	// Create mock WeatherClient
	mockWeatherClient := new(m.MockWeatherClient)
	mockWeatherClient.On("FetchWeatherData", 39.9526, -75.1652).Return(models.WeatherData{}, errors.New("failed to fetch weather data"))

	// Create mock Database (not used in this test)
	mockDB := new(m.MockDatabase)

	// Create the WeatherService
	service := NewWeatherService(mockWeatherClient, mockDB)

	// Call the method
	_, err := service.GetWeatherData(39.9526, -75.1652)

	// Assert the results
	assert.Error(t, err)
	assert.EqualError(t, err, "failed to fetch weather data")
	mockWeatherClient.AssertExpectations(t)
}

// TestStoreWeatherData_Success tests successfully storing weather data
// in the mock database, expecting no errors and the correct snapshot ID.
func TestStoreWeatherData_Success(t *testing.T) {
	// Create mock WeatherClient (not used in this test)
	mockWeatherClient := new(m.MockWeatherClient)

	// Create mock Database
	mockDB := new(m.MockDatabase)
	expectedSnapshotID := 123
	weatherData := models.WeatherData{}
	mockDB.On("StoreWeatherData", weatherData).Return(expectedSnapshotID, nil)

	// Create the WeatherService
	service := NewWeatherService(mockWeatherClient, mockDB)

	// Call the method
	actualSnapshotID, err := service.StoreWeatherData(weatherData)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedSnapshotID, actualSnapshotID)
	mockDB.AssertExpectations(t)
}

// TestStoreWeatherData_Error tests the scenario where storing weather data
// in the mock database results in an error.
func TestStoreWeatherData_Error(t *testing.T) {
	// Create mock WeatherClient (not used in this test)
	mockWeatherClient := new(m.MockWeatherClient)

	// Create mock Database
	mockDB := new(m.MockDatabase)
	weatherData := models.WeatherData{}
	mockDB.On("StoreWeatherData", weatherData).Return(0, errors.New("failed to store weather data"))

	// Create the WeatherService
	service := NewWeatherService(mockWeatherClient, mockDB)

	// Call the method
	_, err := service.StoreWeatherData(weatherData)

	// Assert the results
	assert.Error(t, err)
	assert.EqualError(t, err, "failed to store weather data")
	mockDB.AssertExpectations(t)
}
