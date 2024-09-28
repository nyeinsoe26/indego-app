package db

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nyeinsoe26/indego-app/internal/app/models"
	"github.com/stretchr/testify/assert"
)

// Mock setup for PostgresDB
func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *PostgresDB) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	postgresDB := &PostgresDB{conn: sqlDB}

	return sqlDB, mock, postgresDB
}

// Test for StoreIndegoData success
func TestStoreIndegoData_Success(t *testing.T) {
	_, mock, postgresDB := setupMockDB(t)
	defer postgresDB.Close()

	indegoData := models.IndegoData{
		LastUpdated: time.Now(),
		Features:    []models.StationFeature{}, // Simplified for this test
	}

	indegoJSON, _ := json.Marshal(indegoData)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO indego_snapshots").
		WithArgs(indegoData.LastUpdated, indegoJSON).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	indegoSnapshotID, err := postgresDB.StoreIndegoData(indegoData)
	assert.NoError(t, err)
	assert.Equal(t, 1, indegoSnapshotID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Test for StoreWeatherData success
func TestStoreWeatherData_Success(t *testing.T) {
	_, mock, postgresDB := setupMockDB(t)
	defer postgresDB.Close()

	weatherData := models.WeatherData{
		Main: struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		}{
			Temp:      24.5,
			FeelsLike: 25.0,
			TempMin:   22.0,
			TempMax:   26.0,
			Pressure:  1015,
			Humidity:  60,
		},
	}

	weatherJSON, _ := json.Marshal(weatherData)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO weather_snapshots").
		WithArgs(sqlmock.AnyArg(), weatherJSON).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	mock.ExpectCommit()

	weatherSnapshotID, err := postgresDB.StoreWeatherData(weatherData)
	assert.NoError(t, err)
	assert.Equal(t, 2, weatherSnapshotID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Test for StoreSnapshotLink success
func TestStoreSnapshotLink_Success(t *testing.T) {
	_, mock, postgresDB := setupMockDB(t)
	defer postgresDB.Close()

	timestamp := time.Now()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO snapshots").
		WithArgs(timestamp, 1, 2).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := postgresDB.StoreSnapshotLink(1, 2, timestamp)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFetchSnapshot_Success(t *testing.T) {
	_, mock, postgresDB := setupMockDB(t)
	defer postgresDB.Close()

	indegoData := models.IndegoData{
		LastUpdated: time.Now(),
		Features:    []models.StationFeature{}, // Simplified for this test
	}
	weatherData := models.WeatherData{
		Main: struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		}{
			Temp:      24.5,
			FeelsLike: 25.0,
			TempMin:   22.0,
			TempMax:   26.0,
			Pressure:  1015,
			Humidity:  60,
		},
	}

	indegoJSON, _ := json.Marshal(indegoData)
	weatherJSON, _ := json.Marshal(weatherData)

	mock.ExpectQuery("SELECT i.data, w.data FROM snapshots s").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"i.data", "w.data"}).
			AddRow(indegoJSON, weatherJSON))

	indegoResult, weatherResult, err := postgresDB.FetchSnapshot(time.Now())
	assert.NoError(t, err)

	// Compare the main fields
	assert.Equal(t, indegoData.Features, indegoResult.Features)
	assert.Equal(t, weatherData.Main.Temp, weatherResult.Main.Temp)

	// Compare times with tolerance
	assert.WithinDuration(t, indegoData.LastUpdated, indegoResult.LastUpdated, time.Second)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Test for FetchSnapshot failure (no rows)
func TestFetchSnapshot_NoRows(t *testing.T) {
	_, mock, postgresDB := setupMockDB(t)
	defer postgresDB.Close()

	mock.ExpectQuery("SELECT i.data, w.data FROM snapshots s").
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(sql.ErrNoRows)

	_, _, err := postgresDB.FetchSnapshot(time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no snapshot found")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Test for FetchSnapshot failure (invalid JSON)
func TestFetchSnapshot_InvalidJSON(t *testing.T) {
	_, mock, postgresDB := setupMockDB(t)
	defer postgresDB.Close()

	mock.ExpectQuery("SELECT i.data, w.data FROM snapshots s").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"i.data", "w.data"}).
			AddRow("{invalid_json}", "{invalid_json}"))

	_, _, err := postgresDB.FetchSnapshot(time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal Indego data")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
