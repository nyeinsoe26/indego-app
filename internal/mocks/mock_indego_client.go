package mocks

import (
	"github.com/nyeinsoe26/indego-app/internal/app/models"
	"github.com/stretchr/testify/mock"
)

// MockIndegoClient is a mock implementation of the IndegoClient interface.
type MockIndegoClient struct {
	mock.Mock
}

func (m *MockIndegoClient) FetchIndegoData() (models.IndegoData, error) {
	args := m.Called()
	return args.Get(0).(models.IndegoData), args.Error(1)
}
