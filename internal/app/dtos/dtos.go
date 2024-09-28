package dtos

import "github.com/nyeinsoe26/indego-app/internal/app/models"

// FetchIndegoWeatherResponse represents the success response for fetching and storing Indego and weather data.
type FetchIndegoWeatherResponse struct {
	Message string `json:"message" example:"Data stored successfully"`
}

// StationSnapshotResponse represents the response for station and weather snapshot.
type StationSnapshotResponse struct {
	At       string             `json:"at" example:"2019-09-01T10:00:00Z"`
	Stations models.IndegoData  `json:"stations"`
	Weather  models.WeatherData `json:"weather"`
}

// SpecificStationSnapshotResponse represents the response for a specific station's snapshot.
type SpecificStationSnapshotResponse struct {
	At      string                `json:"at" example:"2019-09-01T10:00:00Z"`
	Station models.StationFeature `json:"station"`
	Weather models.WeatherData    `json:"weather"`
}

// ErrorResponse represents a standard error response.
type ErrorResponse struct {
	Error string `json:"error" example:"Internal Server Error"`
}
