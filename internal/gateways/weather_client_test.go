package gateways

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nyeinsoe26/indego-app/config"
	"github.com/nyeinsoe26/indego-app/internal/app/models"
)

// Test the Weather client with a mocked API response
func TestFetchWeatherData(t *testing.T) {
	// Mock response from the OpenWeather API
	mockResponse := models.WeatherData{
		Main: struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		}{
			Temp:      25.5,
			FeelsLike: 24.0,
			TempMin:   22.5,
			TempMax:   28.5,
			Pressure:  1015,
			Humidity:  60,
		},
		Weather: []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		}{
			{Main: "Clear", Description: "clear sky", Icon: "01d"},
		},
	}

	// Create a test server to simulate the API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Override the Weather base URL in config for testing
	config.AppConfig.Weather.BaseURL = server.URL

	// Test the client
	weatherClient := NewWeatherClient()
	data, err := weatherClient.FetchWeatherData(44.34, 10.99)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the fetched data
	if data.Main.Temp != 25.5 {
		t.Errorf("Expected temperature 25.5, got %f", data.Main.Temp)
	}
	if data.Weather[0].Description != "clear sky" {
		t.Errorf("Expected weather 'clear sky', got %s", data.Weather[0].Description)
	}
}

// Test error handling when API is down
func TestFetchWeatherData_Failure(t *testing.T) {
	// Simulate a server error with a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	// Override the Weather base URL in config for testing
	config.AppConfig.Weather.BaseURL = server.URL

	// Test the client
	weatherClient := NewWeatherClient()
	_, err := weatherClient.FetchWeatherData(44.34, 10.99)
	if err == nil {
		t.Fatalf("Expected error, got none")
	}
}
