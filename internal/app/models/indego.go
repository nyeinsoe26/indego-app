package models

import "time"

// IndegoData defines the structure for the entire Indego API response.
type IndegoData struct {
	LastUpdated time.Time        `json:"last_updated"`
	Features    []StationFeature `json:"features"`
}

// StationFeature defines the station and its related geometry and properties.
type StationFeature struct {
	Geometry   Geometry          `json:"geometry"`
	Properties StationProperties `json:"properties"`
	Type       string            `json:"type"`
}

// Geometry defines the geographic coordinates of the station.
type Geometry struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

// StationProperties defines the details of each station.
type StationProperties struct {
	ID                     int        `json:"id"`
	Name                   string     `json:"name"`
	Coordinates            []float64  `json:"coordinates"`
	TotalDocks             int        `json:"totalDocks"`
	DocksAvailable         int        `json:"docksAvailable"`
	BikesAvailable         int        `json:"bikesAvailable"`
	ClassicBikesAvailable  int        `json:"classicBikesAvailable"`
	SmartBikesAvailable    int        `json:"smartBikesAvailable"`
	ElectricBikesAvailable int        `json:"electricBikesAvailable"`
	RewardBikesAvailable   int        `json:"rewardBikesAvailable"`
	RewardDocksAvailable   int        `json:"rewardDocksAvailable"`
	KioskStatus            string     `json:"kioskStatus"`
	KioskPublicStatus      string     `json:"kioskPublicStatus"`
	KioskConnectionStatus  string     `json:"kioskConnectionStatus"`
	KioskType              int        `json:"kioskType"`
	AddressStreet          string     `json:"addressStreet"`
	AddressCity            string     `json:"addressCity"`
	AddressState           string     `json:"addressState"`
	AddressZipCode         string     `json:"addressZipCode"`
	Bikes                  []Bike     `json:"bikes"`
	CloseTime              *time.Time `json:"closeTime"`
	EventEnd               *time.Time `json:"eventEnd"`
	EventStart             *time.Time `json:"eventStart"`
	IsEventBased           bool       `json:"isEventBased"`
	IsVirtual              bool       `json:"isVirtual"`
	KioskID                int        `json:"kioskId"`
	Notes                  *string    `json:"notes"`
	OpenTime               *time.Time `json:"openTime"`
	PublicText             string     `json:"publicText"`
	TimeZone               *string    `json:"timeZone"`
	TrikesAvailable        int        `json:"trikesAvailable"`
	Latitude               float64    `json:"latitude"`
	Longitude              float64    `json:"longitude"`
}

// Bike defines the structure for the available bikes at the station.
type Bike struct {
	DockNumber  int  `json:"dockNumber"`
	IsElectric  bool `json:"isElectric"`
	IsAvailable bool `json:"isAvailable"`
	Battery     *int `json:"battery"` // Battery can be null
}
