package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/nyeinsoe26/indego-app/internal/app/models"
	"github.com/stretchr/testify/mock"
)

// MockWeatherService is a mock implementation of the WeatherService interface.
type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetWeatherData(lat, lon float64) (models.WeatherData, error) {
	args := m.Called(lat, lon)
	return args.Get(0).(models.WeatherData), args.Error(1)
}

func (m *MockWeatherService) StoreWeatherData(data models.WeatherData, timestamp time.Time) (uuid.UUID, error) {
	args := m.Called(data, timestamp)
	return args.Get(0).(uuid.UUID), args.Error(1)
}
