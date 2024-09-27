package services

import (
	"github.com/nyeinsoe26/indego-app/internal/app/db"
	"github.com/nyeinsoe26/indego-app/internal/app/models"
	"github.com/nyeinsoe26/indego-app/internal/gateways"
)

// WeatherService defines the interface for Weather-related business logic.
type WeatherService interface {
	GetWeatherData(lat, lon float64) (models.WeatherData, error)
	StoreWeatherData(data models.WeatherData) error
}

type weatherServiceImpl struct {
	WeatherClient gateways.WeatherClient
	DB            db.Database
}

// NewWeatherService returns a new WeatherService implementation.
func NewWeatherService(weatherClient gateways.WeatherClient, db db.Database) WeatherService {
	return &weatherServiceImpl{
		WeatherClient: weatherClient,
		DB:            db,
	}
}

// GetWeatherData fetches data from the WeatherClient and returns it.
func (s *weatherServiceImpl) GetWeatherData(lat, lon float64) (models.WeatherData, error) {
	return s.WeatherClient.FetchWeatherData(lat, lon)
}

// StoreWeatherData stores the fetched weather data in the database.
func (s *weatherServiceImpl) StoreWeatherData(data models.WeatherData) error {
	return s.DB.StoreWeatherData(data)
}
