package config

import (
	"os"
	"testing"
)

// Test loading a valid configuration file
func TestLoadConfig_ValidConfig(t *testing.T) {
	cfgPath := "../config.yaml"
	err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("Failed to load valid config: %v", err)
	}

	// Validate server config
	if AppConfig.Server.Port != "3000" {
		t.Errorf("Expected server port 3000, got %s", AppConfig.Server.Port)
	}

	// Validate database config
	if AppConfig.Database.Host != "localhost" {
		t.Errorf("Expected database host 'localhost', got %s", AppConfig.Database.Host)
	}

	// Validate Indego API config
	if AppConfig.Indego.BaseURL != "https://www.rideindego.com/stations/json/" {
		t.Errorf("Expected Indego base URL 'https://www.rideindego.com/stations/json/', got %s", AppConfig.Indego.BaseURL)
	}

	// Validate Weather API config
	if AppConfig.Weather.BaseURL != "https://api.openweathermap.org/data/2.5/weather" {
		t.Errorf("Expected Weather base URL 'https://api.openweathermap.org/data/2.5/weather', got %s", AppConfig.Weather.BaseURL)
	}
}

// Test overriding config values with environment variables
func TestLoadConfig_EnvOverrides(t *testing.T) {
	// Set environment variables to override config values
	os.Setenv("WEATHER_API_KEY", "testapikey")
	os.Setenv("AUTH_TOKEN", "testauthtoken")

	// Reload the config
	cfgPath := "../config.yaml"
	err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check if the environment variables have overridden the config values
	if AppConfig.Weather.APIKey != "testapikey" {
		t.Errorf("Expected overridden weather API key 'testapikey', got %s", AppConfig.Weather.APIKey)
	}

	// Clean up environment variables
	os.Unsetenv("WEATHER_API_KEY")
	os.Unsetenv("AUTH_TOKEN")
}

// Test missing config file handling
func TestLoadConfig_MissingConfigFile(t *testing.T) {
	// Call LoadConfig with a non-existent config file
	err := LoadConfig("nonexistent_config.yaml")
	if err == nil {
		t.Errorf("Expected error due to missing config file, but got none")
	}
}

// Test invalid YAML format in config file using an in-memory string
func TestLoadConfig_InvalidConfigFormat(t *testing.T) {
	// Create an invalid config file for testing
	file, err := os.CreateTemp("", "invalid_config*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file for invalid config: %v", err)
	}
	defer os.Remove(file.Name())

	// Write invalid YAML content to the temp file
	_, err = file.WriteString(`
server:
  port 3000
database:
  host: "localhost"
  - invalid_yaml
`)
	if err != nil {
		t.Fatalf("Failed to write invalid content to temp file: %v", err)
	}
	file.Close()

	// Try loading the invalid config
	err = LoadConfig(file.Name())
	if err == nil {
		t.Errorf("Expected error due to invalid config format, but got none")
	}
}
