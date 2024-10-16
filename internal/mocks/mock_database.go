// Code generated by mockery v2. Do not edit.

package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/nyeinsoe26/indego-app/internal/app/models"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) StoreIndegoData(data models.IndegoData) (uuid.UUID, error) {
	args := m.Called(data)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockDatabase) StoreWeatherData(data models.WeatherData, timestamp time.Time) (uuid.UUID, error) {
	args := m.Called(data, timestamp)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockDatabase) StoreSnapshot(indegoSnapshot models.IndegoData, weatherSnapshot models.WeatherData, timestamp time.Time) (uuid.UUID, error) {
	args := m.Called(indegoSnapshot, weatherSnapshot, timestamp)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockDatabase) FetchSnapshot(at time.Time) (models.IndegoData, models.WeatherData, time.Time, error) {
	args := m.Called(at)
	return args.Get(0).(models.IndegoData), args.Get(1).(models.WeatherData), args.Get(2).(time.Time), args.Error(3)
}
