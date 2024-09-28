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
	StoreIndegoData(data models.IndegoData) (int, error) // Updated to return snapshot ID
	StoreSnapshotLink(indegoSnapshotID, weatherSnapshotID int, timestamp time.Time) error
	GetSnapshot(at time.Time) (models.IndegoData, models.WeatherData, error)
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
func (s *indegoServiceImpl) StoreIndegoData(data models.IndegoData) (int, error) {
	return s.DB.StoreIndegoData(data)
}

// StoreSnapshotLink stores the relationship between Indego and Weather snapshots in the database.
func (s *indegoServiceImpl) StoreSnapshotLink(indegoSnapshotID, weatherSnapshotID int, timestamp time.Time) error {
	return s.DB.StoreSnapshotLink(indegoSnapshotID, weatherSnapshotID, timestamp)
}

// GetSnapshot fetches the first available snapshot of Indego and Weather data at or after the specified time.
func (s *indegoServiceImpl) GetSnapshot(at time.Time) (models.IndegoData, models.WeatherData, error) {
	return s.DB.FetchSnapshot(at)
}
