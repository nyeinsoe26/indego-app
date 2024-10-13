package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/nyeinsoe26/indego-app/internal/app/models"
)

// PostgresDB implements the Database interface for PostgreSQL.
type PostgresDB struct {
	conn *sql.DB
}

// NewPostgresDB returns a new PostgresDB instance with connection pooling configured.
func NewPostgresDB(connStr string) (*PostgresDB, error) {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Verify the connection
	if err = conn.Ping(); err != nil {
		return nil, err
	}

	// Connection Pool settings
	conn.SetMaxOpenConns(25)                 // Max number of open connections to the DB
	conn.SetMaxIdleConns(25)                 // Max number of idle connections in the pool
	conn.SetConnMaxLifetime(5 * time.Minute) // Max lifetime of a connection before being recycled

	return &PostgresDB{conn: conn}, nil
}

// StoreIndegoData stores the Indego data snapshot into PostgreSQL with transaction and rollback.
func (p *PostgresDB) StoreIndegoData(data models.IndegoData) (uuid.UUID, error) {
	// Start a transaction
	tx, err := p.conn.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	// Ensure rollback in case of error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Generate UUID for the Indego snapshot
	indegoSnapshotID := uuid.New()

	// Convert Indego data to JSON
	indegoJSON, err := json.Marshal(data)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to marshal Indego data to JSON: %v", err)
	}

	// Insert the data into the database and return the snapshot UUID
	query := `
		INSERT INTO indego_snapshots (id, timestamp, data) 
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(query, indegoSnapshotID, data.LastUpdated, indegoJSON)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to store Indego data: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return uuid.Nil, err
	}
	return indegoSnapshotID, nil
}

// StoreWeatherData stores the weather data snapshot into PostgreSQL with transaction and rollback.
func (p *PostgresDB) StoreWeatherData(data models.WeatherData, timestamp time.Time) (uuid.UUID, error) {
	// Start a transaction
	tx, err := p.conn.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	// Ensure rollback in case of error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Generate UUID for the Weather snapshot
	weatherSnapshotID := uuid.New()

	// Convert weather data to JSON
	weatherJSON, err := json.Marshal(data)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to marshal weather data to JSON: %v", err)
	}

	// Insert the data into the database and return the snapshot UUID
	query := `
		INSERT INTO weather_snapshots (id, timestamp, data) 
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(query, weatherSnapshotID, timestamp, weatherJSON)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to store weather data: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return uuid.Nil, err
	}
	return weatherSnapshotID, nil
}

// StoreSnapshot inserts both Indego and Weather snapshots and links them in the snapshots table.
// This function ensures atomicity by running the operations within a single transaction.
func (p *PostgresDB) StoreSnapshot(indegoData models.IndegoData, weatherData models.WeatherData, timestamp time.Time) (uuid.UUID, error) {
	// Start a transaction
	tx, err := p.conn.Begin()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Ensure rollback in case of error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. Insert the Indego snapshot with UUID
	indegoSnapshotID := uuid.New()
	indegoJSON, err := json.Marshal(indegoData)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to marshal Indego data to JSON: %v", err)
	}

	indegoInsertQuery := `
		INSERT INTO indego_snapshots (id, timestamp, data)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(indegoInsertQuery, indegoSnapshotID, timestamp, indegoJSON)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert Indego snapshot: %v", err)
	}

	// 2. Insert the Weather snapshot with UUID
	weatherSnapshotID := uuid.New()
	weatherJSON, err := json.Marshal(weatherData)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to marshal Weather data to JSON: %v", err)
	}

	weatherInsertQuery := `
		INSERT INTO weather_snapshots (id, timestamp, data)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(weatherInsertQuery, weatherSnapshotID, timestamp, weatherJSON) // Use Indego timestamp
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert Weather snapshot: %v", err)
	}

	// 3. Link the snapshots by inserting into the `snapshots` table
	snapshotID := uuid.New()
	snapshotInsertQuery := `
		INSERT INTO snapshots (id, timestamp, indego_snapshot_id, weather_snapshot_id)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(snapshotInsertQuery, snapshotID, timestamp, indegoSnapshotID, weatherSnapshotID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert snapshot link: %v", err)
	}

	// 4. Commit the transaction
	err = tx.Commit()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Return the newly created snapshot ID
	return snapshotID, nil
}

// FetchSnapshot fetches the first available snapshot of Indego and Weather data at or after the specified time.
func (p *PostgresDB) FetchSnapshot(at time.Time) (models.IndegoData, models.WeatherData, time.Time, error) {
	query := `
		SELECT s.timestamp, i.data, w.data 
		FROM snapshots s
		JOIN indego_snapshots i ON s.indego_snapshot_id = i.id
		JOIN weather_snapshots w ON s.weather_snapshot_id = w.id
		WHERE s.timestamp >= $1
		ORDER BY s.timestamp ASC
		LIMIT 1
	`

	var indegoJSON, weatherJSON []byte
	var snapshotTime time.Time

	err := p.conn.QueryRow(query, at).Scan(&snapshotTime, &indegoJSON, &weatherJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.IndegoData{}, models.WeatherData{}, time.Time{}, fmt.Errorf("no snapshot found for the given time: %v", at)
		}
		return models.IndegoData{}, models.WeatherData{}, time.Time{}, fmt.Errorf("failed to fetch snapshot: %v", err)
	}

	// Unmarshal the JSON data into the respective models
	var indegoData models.IndegoData
	var weatherData models.WeatherData
	err = json.Unmarshal(indegoJSON, &indegoData)
	if err != nil {
		return models.IndegoData{}, models.WeatherData{}, time.Time{}, fmt.Errorf("failed to unmarshal Indego data: %v", err)
	}

	err = json.Unmarshal(weatherJSON, &weatherData)
	if err != nil {
		return models.IndegoData{}, models.WeatherData{}, time.Time{}, fmt.Errorf("failed to unmarshal weather data: %v", err)
	}

	return indegoData, weatherData, snapshotTime, nil
}

// Close closes the database connection.
func (p *PostgresDB) Close() error {
	return p.conn.Close()
}
