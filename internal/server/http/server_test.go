package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Dendyator/calendar/internal/logger"  //nolint
	"github.com/Dendyator/calendar/internal/storage" //nolint
	"github.com/google/uuid"                         //nolint
	"github.com/stretchr/testify/assert"             //nolint
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) DeleteOldEvents(before time.Time) error {
	args := m.Called(before)
	return args.Error(0)
}

func (m *MockStorage) ListEventsByDay(date time.Time) ([]storage.Event, error) {
	args := m.Called(date)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockStorage) ListEventsByWeek(start time.Time) ([]storage.Event, error) {
	args := m.Called(start)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockStorage) ListEventsByMonth(start time.Time) ([]storage.Event, error) {
	args := m.Called(start)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockStorage) ListEvents() ([]storage.Event, error) {
	args := m.Called()
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockStorage) CreateEvent(event storage.Event) error {
	return m.Called(event).Error(0)
}

func (m *MockStorage) GetEvent(id uuid.UUID) (storage.Event, error) {
	args := m.Called(id)
	return args.Get(0).(storage.Event), args.Error(1)
}

func (m *MockStorage) UpdateEvent(id uuid.UUID, event storage.Event) error {
	return m.Called(id, event).Error(0)
}

func (m *MockStorage) DeleteEvent(id uuid.UUID) error {
	return m.Called(id).Error(0)
}

func TestListEventsHandler_EmptyList(t *testing.T) {
	mockStorage := new(MockStorage)
	logg := logger.New("info")

	NewServer(ServerConfig{Host: "localhost", Port: "8080"}, logg, mockStorage)

	mockStorage.On("ListEvents").Return([]storage.Event{}, nil)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/events", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := listEventsHandler(mockStorage, logg)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var events []storage.Event
	err = json.Unmarshal(rr.Body.Bytes(), &events)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(events))

	mockStorage.AssertExpectations(t)
}

func TestListEventsHandler_WithEvents(t *testing.T) {
	mockStorage := new(MockStorage)
	logg := logger.New("info")

	NewServer(ServerConfig{Host: "localhost", Port: "8080"}, logg, mockStorage)

	event := storage.Event{
		ID:          uuid.New(),
		Title:       "Test Event",
		Description: "This is a test event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	mockStorage.On("ListEvents").Return([]storage.Event{event}, nil)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/events", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := listEventsHandler(mockStorage, logg)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var events []storage.Event
	err = json.Unmarshal(rr.Body.Bytes(), &events)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(events))
	assert.Equal(t, event.ID, events[0].ID)

	mockStorage.AssertExpectations(t)
}

func TestCreateEventHandler(t *testing.T) {
	mockStorage := new(MockStorage)
	logg := logger.New("info")
	fixedTime := time.Date(2024, 11, 11, 13, 4, 0, 0, time.UTC)

	event := storage.Event{
		ID:          uuid.New(),
		Title:       "New Event",
		Description: "This is a new event",
		StartTime:   fixedTime,
		EndTime:     fixedTime.Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	eventJSON, _ := json.Marshal(event)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
		"/events", bytes.NewBuffer(eventJSON))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := createEventHandler(mockStorage, logg)

	mockStorage.On("CreateEvent", event).Return(nil)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestGetEventHandler(t *testing.T) {
	mockStorage := new(MockStorage)
	logg := logger.New("info")

	event := storage.Event{
		ID:          uuid.New(),
		Title:       "Get Event",
		Description: "This is a get event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	mockStorage.On("GetEvent", event.ID).Return(event, nil)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet,
		"/events/"+event.ID.String(), nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := getEventHandler(mockStorage, logg)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var returnedEvent storage.Event
	err = json.Unmarshal(rr.Body.Bytes(), &returnedEvent)
	assert.NoError(t, err)
	assert.Equal(t, event.ID, returnedEvent.ID)
}

func TestUpdateEventHandler(t *testing.T) {
	mockStorage := new(MockStorage)
	logg := logger.New("info")
	fixedTime := time.Date(2024, 11, 11, 13, 4, 0, 0, time.UTC)

	event := storage.Event{
		ID:          uuid.New(),
		Title:       "Updated Event",
		Description: "This is an updated event",
		StartTime:   fixedTime,
		EndTime:     fixedTime.Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	eventJSON, _ := json.Marshal(event)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPut,
		"/events/"+event.ID.String(), bytes.NewBuffer(eventJSON))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := updateEventHandler(mockStorage, logg)

	mockStorage.On("UpdateEvent", event.ID, event).Return(nil)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteEventHandler(t *testing.T) {
	mockStorage := new(MockStorage)
	logg := logger.New("info")

	eventID := uuid.New()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete,
		"/events/"+eventID.String(), nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := deleteEventHandler(mockStorage, logg)

	mockStorage.On("DeleteEvent", eventID).Return(nil)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
