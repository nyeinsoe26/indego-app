package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	_ "github.com/lib/pq"
	"github.com/nyeinsoe26/indego-app/config"
)

const migrationTableCreationQuery = `
CREATE TABLE IF NOT EXISTS migration_history (
    current_version INT NOT NULL,
    execution_time TIMESTAMPTZ NOT NULL
);
`

func main() {
	// Define up and down flags as booleans
	up := flag.Bool("up", false, "Apply the 'up' migration")
	down := flag.Bool("down", false, "Apply the 'down' migration")
	configPath := flag.String("config", "../config.yaml", "Path to the configuration file") // Path to the config file

	flag.Parse()

	// Ensure either up or down is specified
	if !*up && !*down {
		log.Fatalf("Specify either --up or --down to run the migration.")
	}

	// Load the configuration from config.yaml
	config.LoadConfig(*configPath)

	// Connect to the database using config
	connStr := config.GetDatabaseConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Create the migration history table if it doesn't exist
	_, err = db.Exec(migrationTableCreationQuery)
	if err != nil {
		log.Fatalf("Failed to create migration history table: %v", err)
	}

	// Get the current version from the migration history
	currentVersion, err := getCurrentMigrationVersion(db)
	if err != nil {
		log.Fatalf("Failed to get current migration version: %v", err)
	}

	// Get sorted migration files
	migrationFiles, err := getSortedMigrationFiles(".")
	if err != nil {
		log.Fatalf("Failed to get migration files: %v", err)
	}

	// Apply migration based on the flags
	if *up {
		applyUpMigrations(db, migrationFiles, currentVersion)
	} else if *down {
		applyDownMigration(db, migrationFiles, currentVersion)
	}
}

// getCurrentMigrationVersion retrieves the latest migration version from the migration history table
func getCurrentMigrationVersion(db *sql.DB) (int, error) {
	var version int
	err := db.QueryRow(`SELECT current_version FROM migration_history ORDER BY current_version DESC LIMIT 1`).Scan(&version)
	if err != nil {
		// If no rows are returned, we assume it's at version 0 (initial state)
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return version, nil
}

// getSortedMigrationFiles retrieves and sorts the migration SQL files in the correct order
func getSortedMigrationFiles(directory string) ([]string, error) {
	var migrationFiles []string
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	// Regex to match migration files like '001.init.up.sql'
	re := regexp.MustCompile(`^(\d{3}).*\.sql$`)
	for _, file := range files {
		if !file.IsDir() && re.MatchString(file.Name()) {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

// applyUpMigrations applies all migrations that are newer than the current version
func applyUpMigrations(db *sql.DB, migrationFiles []string, currentVersion int) {
	for _, file := range migrationFiles {
		version, err := getVersionFromFile(file)
		if err != nil {
			log.Fatalf("Failed to parse migration file: %v", err)
		}

		// Skip already applied migrations
		if version > currentVersion {
			fmt.Printf("Applying migration: %s\n", file)
			if err := runMigration(db, file, "up"); err != nil {
				log.Fatalf("Failed to apply migration %s: %v", file, err)
			}
			if err := recordMigration(db, version); err != nil {
				log.Fatalf("Failed to record migration: %v", err)
			}
		}
	}
}

// applyDownMigration rolls back the latest migration
func applyDownMigration(db *sql.DB, migrationFiles []string, currentVersion int) {
	if currentVersion == 0 {
		log.Println("No migrations to roll back")
		return
	}

	// Find the migration file corresponding to the current version
	for _, file := range migrationFiles {
		version, err := getVersionFromFile(file)
		if err != nil {
			log.Fatalf("Failed to parse migration file: %v", err)
		}

		if version == currentVersion {
			fmt.Printf("Rolling back migration: %s\n", file)
			if err := runMigration(db, file, "down"); err != nil {
				log.Fatalf("Failed to roll back migration %s: %v", file, err)
			}
			if err := removeMigrationRecord(db, currentVersion); err != nil {
				log.Fatalf("Failed to remove migration record: %v", err)
			}
			break
		}
	}
}

// getVersionFromFile extracts the migration version from the filename (e.g., 001.init.up.sql -> 001)
func getVersionFromFile(fileName string) (int, error) {
	re := regexp.MustCompile(`^(\d{3})`)
	matches := re.FindStringSubmatch(fileName)
	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid migration file name: %s", fileName)
	}
	var version int
	fmt.Sscanf(matches[1], "%d", &version)
	return version, nil
}

// runMigration executes the up or down SQL portion of the migration file
func runMigration(db *sql.DB, filePath string, direction string) error {
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file '%s': %v", filePath, err)
	}

	sqlScript := string(sqlBytes)

	// Split the SQL file into up and down sections
	sections := strings.Split(sqlScript, "--- down ---")
	var queries string
	if direction == "up" {
		queries = sections[0] // First part is for "up"
	} else if len(sections) > 1 {
		queries = sections[1] // Second part is for "down"
	} else {
		return fmt.Errorf("no down section found in migration file '%s'", filePath)
	}

	// Execute the queries
	_, err = db.Exec(queries)
	if err != nil {
		return fmt.Errorf("failed to execute migration: %v", err)
	}

	return nil
}

// recordMigration inserts a record of the applied migration into the migration history
func recordMigration(db *sql.DB, version int) error {
	_, err := db.Exec(`INSERT INTO migration_history (current_version, execution_time) VALUES ($1, NOW())`, version)
	return err
}

// removeMigrationRecord deletes the migration record after a rollback
func removeMigrationRecord(db *sql.DB, version int) error {
	_, err := db.Exec(`DELETE FROM migration_history WHERE current_version = $1`, version)
	return err
}
