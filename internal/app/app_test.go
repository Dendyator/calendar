package app_test

import (
	"context"
	"testing"

	"github.com/Dendyator/calendar/internal/app" //nolint
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLogger struct{}

func (m *MockLogger) Info(_ string)  {}
func (m *MockLogger) Error(_ string) {}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) CreateEvent(ctx context.Context, id, title string) error {
	args := m.Called(ctx, id, title)
	return args.Error(0)
}

func TestCreateEvent(t *testing.T) {
	mockLogger := &MockLogger{}
	mockStorage := new(MockStorage)

	mockStorage.On("CreateEvent", mock.Anything, "1", "Test Event").Return(nil)

	appInstance := app.New(mockLogger, mockStorage)

	err := appInstance.CreateEvent(context.Background(), "1", "Test Event")

	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}
