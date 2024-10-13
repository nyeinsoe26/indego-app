package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config defines the structure of the configuration.
type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Database struct {
		Host     string `mapstructure:"host"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		Port     string `mapstructure:"port"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`

	Indego struct {
		BaseURL string `mapstructure:"base_url"`
	} `mapstructure:"indego"`

	Weather struct {
		BaseURL string `mapstructure:"base_url"`
		APIKey  string `mapstructure:"api_key"`
	} `mapstructure:"weather"`
}

var AppConfig Config

// LoadConfig loads the configuration from a file and allows env overrides.
func LoadConfig(cfgPath string) error {
	viper.SetConfigFile(cfgPath)
	viper.AutomaticEnv()

	// server.port from config yaml can be replaced by SERVER_PORT from env
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Override Weather API key from environment variable
	if apiKey := os.Getenv("WEATHER_API_KEY"); apiKey != "" {
		AppConfig.Weather.APIKey = apiKey
	}

	return nil
}

// GetDatabaseConnectionString returns the PostgreSQL connection string.
func GetDatabaseConnectionString() string {
	return "postgres://" + AppConfig.Database.User + ":" + AppConfig.Database.Password + "@" +
		AppConfig.Database.Host + ":" + AppConfig.Database.Port + "/" + AppConfig.Database.Name +
		"?sslmode=" + AppConfig.Database.SSLMode
}
