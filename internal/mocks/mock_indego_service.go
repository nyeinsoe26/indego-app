package mocks

import (
	"time"

	"github.com/nyeinsoe26/indego-app/internal/app/models"
	"github.com/stretchr/testify/mock"
)

// MockIndegoService is a mock implementation of the IndegoService interface.
type MockIndegoService struct {
	mock.Mock
}

func (m *MockIndegoService) GetIndegoData() (models.IndegoData, error) {
	args := m.Called()
	return args.Get(0).(models.IndegoData), args.Error(1)
}

func (m *MockIndegoService) StoreIndegoData(data models.IndegoData) (int, error) {
	args := m.Called(data)
	return args.Int(0), args.Error(1)
}

func (m *MockIndegoService) GetSnapshot(at time.Time) (models.IndegoData, models.WeatherData, error) {
	args := m.Called(at)
	return args.Get(0).(models.IndegoData), args.Get(1).(models.WeatherData), args.Error(2)
}

func (m *MockIndegoService) StoreSnapshotLink(indegoID, weatherID int, timestamp time.Time) error {
	args := m.Called(indegoID, weatherID, timestamp)
	return args.Error(0)
}
