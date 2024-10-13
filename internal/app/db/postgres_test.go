package db

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
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
	mock.ExpectExec("INSERT INTO indego_snapshots").
		WithArgs(sqlmock.AnyArg(), indegoData.LastUpdated, indegoJSON).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Return success with any result
	mock.ExpectCommit()

	returnedID, err := postgresDB.StoreIndegoData(indegoData)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, returnedID) // assert valid uuid is returned

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
	timestamp := time.Now()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO weather_snapshots").
		WithArgs(sqlmock.AnyArg(), timestamp, weatherJSON).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	returnedID, err := postgresDB.StoreWeatherData(weatherData, timestamp)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, returnedID) // assert a valid uuid is returned

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Test for StoreSnapshot success
func TestStoreSnapshot_Success(t *testing.T) {
	_, mock, postgresDB := setupMockDB(t)
	defer postgresDB.Close()

	indegoData := models.IndegoData{
		LastUpdated: time.Now(),
		Features:    []models.StationFeature{},
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
	timestamp := time.Now()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO indego_snapshots").
		WithArgs(sqlmock.AnyArg(), timestamp, indegoJSON).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO weather_snapshots").
		WithArgs(sqlmock.AnyArg(), timestamp, weatherJSON).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO snapshots").
		WithArgs(sqlmock.AnyArg(), timestamp, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	returnedID, err := postgresDB.StoreSnapshot(indegoData, weatherData, timestamp)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, returnedID) // assert that a valid uuid is returned

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Test for FetchSnapshot success
func TestFetchSnapshot_Success(t *testing.T) {
	_, mock, postgresDB := setupMockDB(t)
	defer postgresDB.Close()

	indegoData := models.IndegoData{
		LastUpdated: time.Now(),
		Features:    []models.StationFeature{},
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
	snapshotTime := time.Now()

	mock.ExpectQuery("SELECT s.timestamp, i.data, w.data FROM snapshots s").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"s.timestamp", "i.data", "w.data"}).
			AddRow(snapshotTime, indegoJSON, weatherJSON))

	indegoResult, weatherResult, returnedSnapshotTime, err := postgresDB.FetchSnapshot(time.Now())
	assert.NoError(t, err)

	// Compare the main fields
	assert.Equal(t, indegoData.Features, indegoResult.Features)
	assert.Equal(t, weatherData.Main.Temp, weatherResult.Main.Temp)
	assert.WithinDuration(t, snapshotTime, returnedSnapshotTime, time.Second)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Test for FetchSnapshot failure (no rows)
func TestFetchSnapshot_NoRows(t *testing.T) {
	_, mock, postgresDB := setupMockDB(t)
	defer postgresDB.Close()

	mock.ExpectQuery("SELECT s.timestamp, i.data, w.data FROM snapshots s").
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(sql.ErrNoRows)

	_, _, _, err := postgresDB.FetchSnapshot(time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no snapshot found")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Test for FetchSnapshot failure (invalid JSON)
func TestFetchSnapshot_InvalidJSON(t *testing.T) {
	_, mock, postgresDB := setupMockDB(t)
	defer postgresDB.Close()

	mock.ExpectQuery("SELECT s.timestamp, i.data, w.data FROM snapshots s").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"s.timestamp", "i.data", "w.data"}).
			AddRow(time.Now(), "{invalid_json}", "{invalid_json}"))

	_, _, _, err := postgresDB.FetchSnapshot(time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal Indego data")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
