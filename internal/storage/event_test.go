package storage

import (
	"testing"
	"time"

	"github.com/google/uuid"              //nolint
	"github.com/stretchr/testify/require" //nolint
)

type MockStorage struct{}

func (m *MockStorage) CreateEvent(_ Event) error {
	return nil
}

func (m *MockStorage) UpdateEvent(_ uuid.UUID, _ Event) error {
	return nil
}

func (m *MockStorage) DeleteEvent(_ uuid.UUID) error {
	return nil
}

func (m *MockStorage) GetEvent(_ uuid.UUID) (Event, error) {
	return Event{}, nil
}

func (m *MockStorage) ListEvents() ([]Event, error) {
	return []Event{}, nil
}

func (m *MockStorage) ListEventsByDay(_ time.Time) ([]Event, error) {
	return []Event{}, nil
}

func (m *MockStorage) ListEventsByWeek(_ time.Time) ([]Event, error) {
	return []Event{}, nil
}

func (m *MockStorage) ListEventsByMonth(_ time.Time) ([]Event, error) {
	return []Event{}, nil
}

func (m *MockStorage) DeleteOldEvents(_ time.Time) error {
	return nil
}

func TestCreateEvent(t *testing.T) {
	mock := &MockStorage{}

	event := Event{
		ID:          uuid.New(),
		Title:       "Test Event",
		Description: "This is a test event.",
		StartTime:   time.Now().Add(time.Hour),
		EndTime:     time.Now().Add(2 * time.Hour),
		UserID:      uuid.New(),
	}

	err := mock.CreateEvent(event)
	require.NoError(t, err, "Event created")
}
