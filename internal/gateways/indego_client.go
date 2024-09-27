package gateways

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nyeinsoe26/indego-app/config"
	"github.com/nyeinsoe26/indego-app/internal/app/models"
)

// IndegoClient defines the interface for interacting with the Indego API.
type IndegoClient interface {
	FetchIndegoData() (models.IndegoData, error)
}

type indegoClientImpl struct{}

// NewIndegoClient returns an implementation of the IndegoClient interface.
func NewIndegoClient() IndegoClient {
	return &indegoClientImpl{}
}

// FetchIndegoData fetches data from the Indego API.
func (c *indegoClientImpl) FetchIndegoData() (models.IndegoData, error) {
	resp, err := http.Get(config.AppConfig.Indego.BaseURL)
	if err != nil {
		return models.IndegoData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.IndegoData{}, errors.New("failed to fetch Indego data")
	}

	var data models.IndegoData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return models.IndegoData{}, err
	}

	return data, nil
}
