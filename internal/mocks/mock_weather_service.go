package mocks

import (
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

func (m *MockWeatherService) StoreWeatherData(data models.WeatherData) (int, error) {
	args := m.Called(data)
	return args.Int(0), args.Error(1)
}
