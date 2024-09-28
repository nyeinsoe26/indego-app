package mocks

import (
	"github.com/nyeinsoe26/indego-app/internal/app/models"
	"github.com/stretchr/testify/mock"
)

// MockWeatherClient is a mock implementation of the WeatherClient interface.
type MockWeatherClient struct {
	mock.Mock
}

func (m *MockWeatherClient) FetchWeatherData(lat, lon float64) (models.WeatherData, error) {
	args := m.Called(lat, lon)
	return args.Get(0).(models.WeatherData), args.Error(1)
}
