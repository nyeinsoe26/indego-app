package db

import (
	"database/sql"
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

// StoreIndegoData saves the Indego data snapshot into PostgreSQL.
func (db *PostgresDB) StoreIndegoData(data models.IndegoData) error {
	// Start a transaction
	tx, err := db.conn.Begin()
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
		INSERT INTO indego_snapshots (timestamp, data)
		VALUES ($1, $2)
	`
	_, err = tx.Exec(query, time.Now(), data)
	if err != nil {
		return err
	}

	// Commit the transaction if successful
	err = tx.Commit()
	return err
}

// StoreWeatherData saves the weather data snapshot into PostgreSQL.
func (db *PostgresDB) StoreWeatherData(data models.WeatherData) error {
	// Start a transaction
	tx, err := db.conn.Begin()
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
		INSERT INTO weather_snapshots (timestamp, data)
		VALUES ($1, $2)
	`
	_, err = tx.Exec(query, time.Now(), data)
	if err != nil {
		return err
	}

	// Commit the transaction if successful
	err = tx.Commit()
	return err
}

// FetchSnapshot retrieves the first snapshot of data on or after the specified time.
func (db *PostgresDB) FetchSnapshot(at time.Time) (models.IndegoData, models.WeatherData, error) {
	// Example SQL query to retrieve the snapshot:
	// SELECT indego_data, weather_data FROM snapshots WHERE timestamp >= $1 ORDER BY timestamp ASC LIMIT 1
	var indegoData models.IndegoData
	var weatherData models.WeatherData

	query := `
		SELECT indego_data, weather_data
		FROM snapshots
		WHERE timestamp >= $1
		ORDER BY timestamp ASC
		LIMIT 1
	`

	err := db.conn.QueryRow(query, at).Scan(&indegoData, &weatherData)
	if err != nil {
		return models.IndegoData{}, models.WeatherData{}, err
	}

	return indegoData, weatherData, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	return p.conn.Close()
}
