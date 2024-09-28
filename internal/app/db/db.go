package db

import (
	"time"

	"github.com/nyeinsoe26/indego-app/internal/app/models"
)

// Database defines the interface for storing and retrieving data.
type Database interface {
	StoreIndegoData(data models.IndegoData) (int, error)
	StoreWeatherData(data models.WeatherData) (int, error)
	StoreSnapshotLink(indegoSnapshotID, weatherSnapshotID int, timestamp time.Time) error
	FetchSnapshot(at time.Time) (models.IndegoData, models.WeatherData, error)
}
