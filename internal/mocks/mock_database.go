// Code generated by mockery v2. Do not edit.

package mocks

import (
	"time"

	"github.com/nyeinsoe26/indego-app/internal/app/models"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) StoreIndegoData(data models.IndegoData) (int, error) {
	args := m.Called(data)
	return args.Int(0), args.Error(1)
}

func (m *MockDatabase) StoreWeatherData(data models.WeatherData) (int, error) {
	args := m.Called(data)
	return args.Int(0), args.Error(1)
}

func (m *MockDatabase) StoreSnapshotLink(indegoSnapshotID, weatherSnapshotID int, timestamp time.Time) error {
	args := m.Called(indegoSnapshotID, weatherSnapshotID, timestamp)
	return args.Error(0)
}

func (m *MockDatabase) FetchSnapshot(at time.Time) (models.IndegoData, models.WeatherData, error) {
	args := m.Called(at)
	return args.Get(0).(models.IndegoData), args.Get(1).(models.WeatherData), args.Error(2)
}
