package grpc

import (
	"context"
	"testing"
	"time"

	pb "github.com/Dendyator/calendar/api/pb"        //nolint
	"github.com/Dendyator/calendar/internal/logger"  //nolint
	"github.com/Dendyator/calendar/internal/storage" //nolint
	"github.com/google/uuid"                         //nolint
	"github.com/stretchr/testify/assert"             //nolint
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) CreateEvent(event storage.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockStorage) UpdateEvent(id uuid.UUID, newEvent storage.Event) error {
	args := m.Called(id, newEvent)
	return args.Error(0)
}

func (m *MockStorage) DeleteEvent(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStorage) GetEvent(id uuid.UUID) (storage.Event, error) {
	args := m.Called(id)
	return args.Get(0).(storage.Event), args.Error(1)
}

func (m *MockStorage) ListEvents() ([]storage.Event, error) {
	args := m.Called()
	return args.Get(0).([]storage.Event), args.Error(1)
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

func (m *MockStorage) DeleteOldEvents(before time.Time) error {
	args := m.Called(before)
	return args.Error(0)
}

func TestCreateEvent(t *testing.T) {
	mockStorage := new(MockStorage)
	logg := logger.New("info")

	server := NewGRPCServer(mockStorage, logg)

	event := &pb.Event{
		Title:       "Test Event",
		Description: "This is a test event.",
		StartTime:   time.Now().Unix(),
		EndTime:     time.Now().Add(1 * time.Hour).Unix(),
		UserId:      uuid.New().String(),
	}

	mockStorage.On("CreateEvent", mock.Anything).Return(nil)

	resp, err := server.CreateEvent(context.Background(), &pb.CreateEventRequest{Event: event})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockStorage.AssertExpectations(t)
}

func TestUpdateEvent(t *testing.T) {
	mockStorage := new(MockStorage)
	logg := logger.New("info")

	server := NewGRPCServer(mockStorage, logg)

	eventID := uuid.New().String()
	newEvent := &pb.Event{
		Id:          eventID,
		Title:       "Updated Event",
		Description: "This is an updated test event.",
		StartTime:   time.Now().Unix(),
		EndTime:     time.Now().Add(1 * time.Hour).Unix(),
		UserId:      uuid.New().String(),
	}

	mockStorage.On("UpdateEvent", mock.Anything, mock.Anything).Return(nil)

	resp, err := server.UpdateEvent(context.Background(), &pb.UpdateEventRequest{Id: eventID, Event: newEvent})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockStorage.AssertExpectations(t)
}

func TestDeleteEvent(t *testing.T) {
	mockStorage := new(MockStorage)
	logg := logger.New("info")

	server := NewGRPCServer(mockStorage, logg)

	eventID := uuid.New().String()

	mockStorage.On("DeleteEvent", mock.Anything).Return(nil)

	resp, err := server.DeleteEvent(context.Background(), &pb.DeleteEventRequest{Id: eventID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockStorage.AssertExpectations(t)
}

func TestGetEvent(t *testing.T) {
	mockStorage := new(MockStorage)
	logg := logger.New("info")

	server := NewGRPCServer(mockStorage, logg)

	eventID := uuid.New()
	expectedEvent := storage.Event{
		ID:          eventID,
		Title:       "Test Event",
		Description: "This is a test event.",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	mockStorage.On("GetEvent", eventID).Return(expectedEvent, nil)

	resp, err := server.GetEvent(context.Background(), &pb.GetEventRequest{Id: eventID.String()})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedEvent.ID.String(), resp.Event.Id)
	mockStorage.AssertExpectations(t)
}

func TestListEvents(t *testing.T) {
	mockStorage := new(MockStorage)
	logg := logger.New("info")

	server := NewGRPCServer(mockStorage, logg)

	events := []storage.Event{
		{
			ID:          uuid.New(),
			Title:       "Event 1",
			Description: "First event",
			StartTime:   time.Now(),
			EndTime:     time.Now().Add(1 * time.Hour),
			UserID:      uuid.New(),
		},
		{
			ID:          uuid.New(),
			Title:       "Event 2",
			Description: "Second event",
			StartTime:   time.Now(),
			EndTime:     time.Now().Add(2 * time.Hour),
			UserID:      uuid.New(),
		},
	}

	mockStorage.On("ListEvents").Return(events, nil)

	resp, err := server.ListEvents(context.Background(), &pb.ListEventsRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(events), len(resp.Events))
	mockStorage.AssertExpectations(t)
}

func TestListEventsByDay(t *testing.T) {
	mockStorage := new(MockStorage)
	logg := logger.New("info")

	server := NewGRPCServer(mockStorage, logg)

	date := time.Now().Truncate(time.Second)
	events := []storage.Event{
		{
			ID:          uuid.New(),
			Title:       "Daily Event 1",
			Description: "First daily event",
			StartTime:   date,
			EndTime:     date.Add(1 * time.Hour),
			UserID:      uuid.New(),
		},
	}

	mockStorage.On("ListEventsByDay", date).Return(events, nil)

	resp, err := server.ListEventsByDay(context.Background(), &pb.ListEventsByDayRequest{Date: date.Unix()})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(events), len(resp.Events))
	mockStorage.AssertExpectations(t)
}

func TestListEventsByWeek(t *testing.T) {
	mockStorage := new(MockStorage)
	logg := logger.New("info")

	server := NewGRPCServer(mockStorage, logg)

	start := time.Now().Truncate(time.Second)
	events := []storage.Event{
		{
			ID:          uuid.New(),
			Title:       "Weekly Event 1",
			Description: "First weekly event",
			StartTime:   start,
			EndTime:     start.Add(1 * time.Hour),
			UserID:      uuid.New(),
		},
	}

	mockStorage.On("ListEventsByWeek", start).Return(events, nil)

	resp, err := server.ListEventsByWeek(context.Background(), &pb.ListEventsByWeekRequest{Start: start.Unix()})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(events), len(resp.Events))
	mockStorage.AssertExpectations(t)
}

func TestListEventsByMonth(t *testing.T) {
	mockStorage := new(MockStorage)
	logg := logger.New("info")

	server := NewGRPCServer(mockStorage, logg)

	start := time.Now().Truncate(time.Second)
	events := []storage.Event{
		{
			ID:          uuid.New(),
			Title:       "Monthly Event 1",
			Description: "First monthly event",
			StartTime:   start,
			EndTime:     start.Add(1 * time.Hour),
			UserID:      uuid.New(),
		},
	}

	mockStorage.On("ListEventsByMonth", start).Return(events, nil)

	resp, err := server.ListEventsByMonth(context.Background(), &pb.ListEventsByMonthRequest{Start: start.Unix()})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(events), len(resp.Events))
	mockStorage.AssertExpectations(t)
}
