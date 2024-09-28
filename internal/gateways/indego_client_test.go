package gateways

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nyeinsoe26/indego-app/config"
	"github.com/nyeinsoe26/indego-app/internal/app/models"
)

// TestFetchIndegoData tests the Indego client's ability to fetch data from the Indego API
// using a mocked API response. It verifies that the fetched station data and bike data match
// the expected mock values, ensuring the client correctly handles a successful API response.
func TestFetchIndegoData(t *testing.T) {
	// Mock response from the Indego API
	mockResponse := models.IndegoData{
		LastUpdated: time.Now(),
		Features: []models.StationFeature{
			{
				Geometry: models.Geometry{
					Coordinates: []float64{-75.14403, 39.94733},
					Type:        "Point",
				},
				Properties: models.StationProperties{
					ID:                     3005,
					Name:                   "Welcome Park, NPS",
					Coordinates:            []float64{-75.14403, 39.94733},
					TotalDocks:             13,
					DocksAvailable:         7,
					BikesAvailable:         4,
					ClassicBikesAvailable:  3,
					ElectricBikesAvailable: 1,
					RewardBikesAvailable:   4,
					RewardDocksAvailable:   8,
					KioskStatus:            "FullService",
					AddressStreet:          "191 S. 2nd St.",
					AddressCity:            "Philadelphia",
					AddressState:           "PA",
					AddressZipCode:         "19106",
					Bikes: []models.Bike{
						{DockNumber: 1, IsElectric: true, IsAvailable: false, Battery: new(int)},
						{DockNumber: 6, IsElectric: false, IsAvailable: true, Battery: nil},
					},
				},
				Type: "Feature",
			},
		},
	}

	// Create a test server to simulate the Indego API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Override the Indego base URL in the config for testing
	config.AppConfig.Indego.BaseURL = server.URL

	// Test the client
	indegoClient := NewIndegoClient()
	data, err := indegoClient.FetchIndegoData()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the fetched station data
	if len(data.Features) == 0 {
		t.Fatalf("Expected station data, got none")
	}
	if data.Features[0].Properties.Name != "Welcome Park, NPS" {
		t.Errorf("Expected station name 'Welcome Park, NPS', got %s", data.Features[0].Properties.Name)
	}

	// Validate bike data within the station feature
	if len(data.Features[0].Properties.Bikes) == 0 {
		t.Fatalf("Expected bike data, got none")
	}
	if data.Features[0].Properties.Bikes[0].DockNumber != 1 {
		t.Errorf("Expected bike dock number 1, got %d", data.Features[0].Properties.Bikes[0].DockNumber)
	}
}

// TestFetchIndegoData_Failure tests the error handling in the Indego client when the API returns
// a failure response, such as a 500 Internal Server Error. It verifies that the client correctly
// handles and returns the error when the API is unavailable or fails.
func TestFetchIndegoData_Failure(t *testing.T) {
	// Simulate a server error with a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	// Override the Indego base URL in the config for testing
	config.AppConfig.Indego.BaseURL = server.URL

	// Test the client
	indegoClient := NewIndegoClient()
	_, err := indegoClient.FetchIndegoData()
	if err == nil {
		t.Fatalf("Expected error, got none")
	}
}
