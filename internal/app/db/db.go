package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/nyeinsoe26/indego-app/internal/app/models"
)

// Database defines the interface for storing and retrieving data.
type Database interface {
	StoreIndegoData(data models.IndegoData) (uuid.UUID, error)
	StoreWeatherData(data models.WeatherData, timestamp time.Time) (uuid.UUID, error)
	StoreSnapshot(indegoSnapshot models.IndegoData, weatherSnapshot models.WeatherData, timestamp time.Time) (uuid.UUID, error)
	FetchSnapshot(at time.Time) (models.IndegoData, models.WeatherData, time.Time, error)
}
