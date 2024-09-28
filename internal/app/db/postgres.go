package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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
func (p *PostgresDB) StoreIndegoData(data models.IndegoData) (int, error) {
	// Start a transaction
	tx, err := p.conn.Begin()
	if err != nil {
		return 0, err
	}

	// Ensure rollback in case of error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		INSERT INTO indego_snapshots (timestamp, data) 
		VALUES ($1, $2)
		RETURNING id
	`

	// Convert Indego data to JSON
	indegoJSON, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal Indego data to JSON: %v", err)
	}

	// Insert the data into the database and return the snapshot ID
	var indegoSnapshotID int
	err = tx.QueryRow(query, data.LastUpdated, indegoJSON).Scan(&indegoSnapshotID)
	if err != nil {
		return 0, fmt.Errorf("failed to store Indego data: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	return indegoSnapshotID, err
}

// StoreWeatherData stores the weather data snapshot into PostgreSQL with transaction and rollback.
func (p *PostgresDB) StoreWeatherData(data models.WeatherData) (int, error) {
	// Start a transaction
	tx, err := p.conn.Begin()
	if err != nil {
		return 0, err
	}

	// Ensure rollback in case of error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		INSERT INTO weather_snapshots (timestamp, data) 
		VALUES ($1, $2)
		RETURNING id
	`

	// Convert weather data to JSON
	weatherJSON, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal weather data to JSON: %v", err)
	}

	// Insert the data into the database and return the snapshot ID
	var weatherSnapshotID int
	err = tx.QueryRow(query, time.Now().UTC(), weatherJSON).Scan(&weatherSnapshotID)
	if err != nil {
		return 0, fmt.Errorf("failed to store weather data: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	return weatherSnapshotID, err
}

// StoreSnapshotLink stores the relationship between Indego and weather snapshots with transaction and rollback.
func (p *PostgresDB) StoreSnapshotLink(indegoSnapshotID, weatherSnapshotID int, timestamp time.Time) error {
	// Start a transaction
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}

	// Ensure rollback in case of error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		INSERT INTO snapshots (timestamp, indego_snapshot_id, weather_snapshot_id) 
		VALUES ($1, $2, $3)
	`

	_, err = tx.Exec(query, timestamp, indegoSnapshotID, weatherSnapshotID)
	if err != nil {
		return fmt.Errorf("failed to store snapshot link: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	return err
}

// FetchSnapshot fetches the first available snapshot of Indego and Weather data at or after the specified time.
func (p *PostgresDB) FetchSnapshot(at time.Time) (models.IndegoData, models.WeatherData, error) {
	query := `
		SELECT i.data, w.data 
		FROM snapshots s
		JOIN indego_snapshots i ON s.indego_snapshot_id = i.id
		JOIN weather_snapshots w ON s.weather_snapshot_id = w.id
		WHERE s.timestamp >= $1
		ORDER BY s.timestamp ASC
		LIMIT 1
	`

	var indegoJSON, weatherJSON []byte
	err := p.conn.QueryRow(query, at).Scan(&indegoJSON, &weatherJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.IndegoData{}, models.WeatherData{}, fmt.Errorf("no snapshot found for the given time: %v", at)
		}
		return models.IndegoData{}, models.WeatherData{}, fmt.Errorf("failed to fetch snapshot: %v", err)
	}

	// Unmarshal the JSON data into the respective models
	var indegoData models.IndegoData
	var weatherData models.WeatherData
	err = json.Unmarshal(indegoJSON, &indegoData)
	if err != nil {
		return models.IndegoData{}, models.WeatherData{}, fmt.Errorf("failed to unmarshal Indego data: %v", err)
	}

	err = json.Unmarshal(weatherJSON, &weatherData)
	if err != nil {
		return models.IndegoData{}, models.WeatherData{}, fmt.Errorf("failed to unmarshal weather data: %v", err)
	}

	return indegoData, weatherData, nil
}

// Close closes the database connection.
func (p *PostgresDB) Close() error {
	return p.conn.Close()
}
