package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/nyeinsoe26/indego-app/internal/app/db"
	"github.com/nyeinsoe26/indego-app/internal/app/models"
	"github.com/nyeinsoe26/indego-app/internal/gateways"
)

// IndegoService defines the interface for Indego-related business logic.
type IndegoService interface {
	GetIndegoData() (models.IndegoData, error)
	StoreIndegoData(data models.IndegoData) (uuid.UUID, error) // Updated to return snapshot ID
	StoreSnapshot(indegoSnapshot models.IndegoData, weatherSnapshot models.WeatherData, timestamp time.Time) (uuid.UUID, error)
	GetSnapshot(at time.Time) (models.IndegoData, models.WeatherData, time.Time, error)
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

// StoreIndegoData stores the fetched Indego data in the database and returns the snapshot ID.
func (s *indegoServiceImpl) StoreIndegoData(data models.IndegoData) (uuid.UUID, error) {
	return s.DB.StoreIndegoData(data)
}

func (s *indegoServiceImpl) StoreSnapshot(indegoSnapshot models.IndegoData, weatherSnapshot models.WeatherData, timestamp time.Time) (uuid.UUID, error) {
	return s.DB.StoreSnapshot(indegoSnapshot, weatherSnapshot, timestamp)
}

// GetSnapshot fetches the first available snapshot of Indego and Weather data at or after the specified time.
func (s *indegoServiceImpl) GetSnapshot(at time.Time) (models.IndegoData, models.WeatherData, time.Time, error) {
	return s.DB.FetchSnapshot(at)
}
