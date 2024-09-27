package gateways

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/nyeinsoe26/indego-app/config"
	"github.com/nyeinsoe26/indego-app/internal/app/models"
)

// WeatherClient defines the interface for interacting with the OpenWeather API.
type WeatherClient interface {
	FetchWeatherData(lat, lon float64) (models.WeatherData, error)
}

type weatherClientImpl struct{}

// NewWeatherClient returns an implementation of the WeatherClient interface.
func NewWeatherClient() WeatherClient {
	return &weatherClientImpl{}
}

// FetchWeatherData fetches weather data from the OpenWeather API.
func (c *weatherClientImpl) FetchWeatherData(lat, lon float64) (models.WeatherData, error) {
	// Use the base URL and API key from the config
	apiURL := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s", config.AppConfig.Weather.BaseURL, lat, lon, config.AppConfig.Weather.APIKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return models.WeatherData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.WeatherData{}, errors.New("failed to fetch weather data")
	}

	var data models.WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return models.WeatherData{}, err
	}

	return data, nil
}
