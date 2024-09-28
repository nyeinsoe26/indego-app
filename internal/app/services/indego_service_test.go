package services

import (
	"errors"
	"testing"
	"time"

	"github.com/nyeinsoe26/indego-app/internal/app/models"
	m "github.com/nyeinsoe26/indego-app/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestGetIndegoData_Success tests the successful retrieval of Indego data
// from the mock IndegoClient, expecting no errors and correct data.
func TestGetIndegoData_Success(t *testing.T) {
	// Mock IndegoClient
	mockClient := new(m.MockIndegoClient)
	expectedData := models.IndegoData{
		LastUpdated: time.Now(),
		Features:    []models.StationFeature{},
	}
	mockClient.On("FetchIndegoData").Return(expectedData, nil)

	// Create the service
	service := NewIndegoService(mockClient, nil)

	// Call the method
	actualData, err := service.GetIndegoData()

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedData, actualData)
	mockClient.AssertExpectations(t)
}

// TestGetIndegoData_Error tests the scenario where fetching Indego data
// results in an error, expecting an error response.
func TestGetIndegoData_Error(t *testing.T) {
	// Mock IndegoClient
	mockClient := new(m.MockIndegoClient)
	mockClient.On("FetchIndegoData").Return(models.IndegoData{}, errors.New("error fetching data"))

	// Create the service
	service := NewIndegoService(mockClient, nil)

	// Call the method
	_, err := service.GetIndegoData()

	// Assert the results
	assert.Error(t, err)
	assert.EqualError(t, err, "error fetching data")
	mockClient.AssertExpectations(t)
}

// TestStoreIndegoData_Success tests successfully storing Indego data
// in the mock database, expecting no errors and the correct snapshot ID.
func TestStoreIndegoData_Success(t *testing.T) {
	// Mock Database
	mockDB := new(m.MockDatabase)
	expectedSnapshotID := 123
	indegoData := models.IndegoData{
		LastUpdated: time.Now(),
		Features:    []models.StationFeature{},
	}
	mockDB.On("StoreIndegoData", indegoData).Return(expectedSnapshotID, nil)

	// Create the service
	service := NewIndegoService(nil, mockDB)

	// Call the method
	actualSnapshotID, err := service.StoreIndegoData(indegoData)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedSnapshotID, actualSnapshotID)
	mockDB.AssertExpectations(t)
}

// TestStoreIndegoData_Error tests the scenario where storing Indego data
// in the mock database results in an error.
func TestStoreIndegoData_Error(t *testing.T) {
	// Mock Database
	mockDB := new(m.MockDatabase)
	indegoData := models.IndegoData{
		LastUpdated: time.Now(),
		Features:    []models.StationFeature{},
	}
	mockDB.On("StoreIndegoData", indegoData).Return(0, errors.New("error storing data"))

	// Create the service
	service := NewIndegoService(nil, mockDB)

	// Call the method
	_, err := service.StoreIndegoData(indegoData)

	// Assert the results
	assert.Error(t, err)
	assert.EqualError(t, err, "error storing data")
	mockDB.AssertExpectations(t)
}

// TestStoreSnapshotLink_Success tests successfully storing the snapshot link
// between Indego and Weather data in the mock database, expecting no errors.
func TestStoreSnapshotLink_Success(t *testing.T) {
	// Mock Database
	mockDB := new(m.MockDatabase)
	mockDB.On("StoreSnapshotLink", 123, 456, mock.Anything).Return(nil)

	// Create the service
	service := NewIndegoService(nil, mockDB)

	// Call the method
	err := service.StoreSnapshotLink(123, 456, time.Now())

	// Assert the results
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

// TestStoreSnapshotLink_Error tests the scenario where storing the snapshot link
// in the mock database results in an error.
func TestStoreSnapshotLink_Error(t *testing.T) {
	// Mock Database
	mockDB := new(m.MockDatabase)
	mockDB.On("StoreSnapshotLink", 123, 456, mock.Anything).Return(errors.New("error storing snapshot link"))

	// Create the service
	service := NewIndegoService(nil, mockDB)

	// Call the method
	err := service.StoreSnapshotLink(123, 456, time.Now())

	// Assert the results
	assert.Error(t, err)
	assert.EqualError(t, err, "error storing snapshot link")
	mockDB.AssertExpectations(t)
}

// TestGetSnapshot_Success tests the successful retrieval of Indego and Weather snapshots
// from the mock database, expecting no errors and correct data.
func TestGetSnapshot_Success(t *testing.T) {
	// Mock Database
	mockDB := new(m.MockDatabase)
	expectedIndegoData := models.IndegoData{
		LastUpdated: time.Now(),
		Features:    []models.StationFeature{},
	}
	expectedWeatherData := models.WeatherData{}
	mockDB.On("FetchSnapshot", mock.Anything).Return(expectedIndegoData, expectedWeatherData, nil)

	// Create the service
	service := NewIndegoService(nil, mockDB)

	// Call the method
	actualIndegoData, actualWeatherData, err := service.GetSnapshot(time.Now())

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedIndegoData, actualIndegoData)
	assert.Equal(t, expectedWeatherData, actualWeatherData)
	mockDB.AssertExpectations(t)
}

// TestGetSnapshot_Error tests the scenario where fetching a snapshot
// from the mock database results in an error.
func TestGetSnapshot_Error(t *testing.T) {
	// Mock Database
	mockDB := new(m.MockDatabase)
	mockDB.On("FetchSnapshot", mock.Anything).Return(models.IndegoData{}, models.WeatherData{}, errors.New("error fetching snapshot"))

	// Create the service
	service := NewIndegoService(nil, mockDB)

	// Call the method
	_, _, err := service.GetSnapshot(time.Now())

	// Assert the results
	assert.Error(t, err)
	assert.EqualError(t, err, "error fetching snapshot")
	mockDB.AssertExpectations(t)
}
