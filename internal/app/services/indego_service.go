package services

import (
	"time"

	"github.com/nyeinsoe26/indego-app/internal/app/db"
	"github.com/nyeinsoe26/indego-app/internal/app/models"
	"github.com/nyeinsoe26/indego-app/internal/gateways"
)

// IndegoService defines the interface for Indego-related business logic.
type IndegoService interface {
	GetIndegoData() (models.IndegoData, error)
	StoreIndegoData(data models.IndegoData) error
	GetSnapshot(at time.Time) (models.IndegoData, models.WeatherData, error) // Added GetSnapshot function
}

type indegoServiceImpl struct {
	IndegoClient gateways.IndegoClient
	DB           db.Database
}

// NewIndegoService returns a new IndegoService implementation.
func NewIndegoService(indegoClient gateways.IndegoClient, db db.Database) IndegoService {
	return &indegoServiceImpl{
		IndegoClient: indegoClient,
		DB:           db,
	}
}

// GetIndegoData fetches data from the IndegoClient and returns it.
func (s *indegoServiceImpl) GetIndegoData() (models.IndegoData, error) {
	return s.IndegoClient.FetchIndegoData()
}

// StoreIndegoData stores the fetched Indego data in the database.
func (s *indegoServiceImpl) StoreIndegoData(data models.IndegoData) error {
	return s.DB.StoreIndegoData(data)
}

// GetSnapshot fetches the first available snapshot of Indego and Weather data at or after the specified time.
func (s *indegoServiceImpl) GetSnapshot(at time.Time) (models.IndegoData, models.WeatherData, error) {
	indegoData, weatherData, err := s.DB.FetchSnapshot(at)
	if err != nil {
		return models.IndegoData{}, models.WeatherData{}, err
	}
	return indegoData, weatherData, nil
}
