package sqlstorage_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage" //nolint
	"github.com/google/uuid"                                         //nolint
)

type MockStorage struct {
	events map[uuid.UUID]storage.Event
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		events: make(map[uuid.UUID]storage.Event),
	}
}

func (m *MockStorage) CreateEvent(event storage.Event) error {
	m.events[event.ID] = event
	return nil
}

func (m *MockStorage) UpdateEvent(id uuid.UUID, newEvent storage.Event) error {
	if _, exists := m.events[id]; !exists {
		return errors.New("event not found")
	}
	m.events[id] = newEvent
	return nil
}

func (m *MockStorage) DeleteEvent(id uuid.UUID) error {
	if _, exists := m.events[id]; !exists {
		return errors.New("event not found")
	}
	delete(m.events, id)
	return nil
}

func (m *MockStorage) GetEvent(id uuid.UUID) (storage.Event, error) {
	event, exists := m.events[id]
	if !exists {
		return storage.Event{}, errors.New("event not found")
	}
	return event, nil
}

func (m *MockStorage) ListEvents() ([]storage.Event, error) {
	events := []storage.Event{}
	for _, event := range m.events {
		events = append(events, event)
	}
	return events, nil
}

func (m *MockStorage) ListEventsByDay(date time.Time) ([]storage.Event, error) {
	start := date.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)
	events := []storage.Event{}
	for _, event := range m.events {
		if event.StartTime.After(start) && event.StartTime.Before(end) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (m *MockStorage) ListEventsByWeek(start time.Time) ([]storage.Event, error) {
	end := start.AddDate(0, 0, 7)
	events := []storage.Event{}
	for _, event := range m.events {
		if event.StartTime.After(start) && event.StartTime.Before(end) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (m *MockStorage) ListEventsByMonth(start time.Time) ([]storage.Event, error) {
	end := start.AddDate(0, 1, 0)
	events := []storage.Event{}
	for _, event := range m.events {
		if event.StartTime.After(start) && event.StartTime.Before(end) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (m *MockStorage) DeleteOldEvents(before time.Time) error {
	for id, event := range m.events {
		if event.EndTime.Before(before) {
			delete(m.events, id)
		}
	}
	return nil
}

func TestCreateEvent(t *testing.T) {
	s := NewMockStorage()
	eventID := uuid.New()

	event := storage.Event{
		ID:          eventID,
		Title:       "Test Event",
		Description: "This is a test event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	err := s.CreateEvent(event)
	if err != nil {
		t.Errorf("error was not expected while creating event: %s", err)
	}

	if len(s.events) != 1 {
		t.Errorf("expected 1 event, got %d", len(s.events))
	}
}

func TestUpdateEvent(t *testing.T) {
	s := NewMockStorage()
	eventID := uuid.New()

	event := storage.Event{
		ID:          eventID,
		Title:       "Test Event",
		Description: "This is a test event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	err := s.CreateEvent(event)
	if err != nil {
		return
	}

	newEvent := storage.Event{
		ID:          eventID,
		Title:       "Updated Event",
		Description: "This is an updated test event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(2 * time.Hour),
		UserID:      event.UserID,
	}

	er := s.UpdateEvent(eventID, newEvent)
	if er != nil {
		t.Errorf("error was not expected while updating event: %s", err)
	}

	updatedEvent, _ := s.GetEvent(eventID)
	if updatedEvent.Title != "Updated Event" {
		t.Errorf("expected event title to be 'Updated Event', got '%s'", updatedEvent.Title)
	}
}

func TestDeleteEvent(t *testing.T) {
	s := NewMockStorage()
	eventID := uuid.New()

	event := storage.Event{
		ID:          eventID,
		Title:       "Test Event",
		Description: "This is a test event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	err := s.CreateEvent(event)
	if err != nil {
		return
	}

	er := s.DeleteEvent(eventID)
	if er != nil {
		t.Errorf("error was not expected while deleting event: %s", err)
	}

	if len(s.events) != 0 {
		t.Errorf("expected 0 events, got %d", len(s.events))
	}
}

func TestDeleteNonExistentEvent(t *testing.T) {
	s := NewMockStorage()
	eventID := uuid.New()

	err := s.DeleteEvent(eventID)
	if err == nil {
		t.Error("expected error while deleting non-existent event, got none")
	}
}

func TestGetEvent(t *testing.T) {
	s := NewMockStorage()
	eventID := uuid.New()

	event := storage.Event{
		ID:          eventID,
		Title:       "Test Event",
		Description: "This is a test event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	err := s.CreateEvent(event)
	if err != nil {
		return
	}

	fetchedEvent, err := s.GetEvent(eventID)
	if err != nil {
		t.Errorf("error was not expected while getting event: %s", err)
	}

	if fetchedEvent.ID != eventID {
		t.Errorf("expected event ID to be '%s', got '%s'", eventID, fetchedEvent.ID)
	}
}

func TestGetNonExistentEvent(t *testing.T) {
	s := NewMockStorage()
	eventID := uuid.New()

	_, err := s.GetEvent(eventID)
	if err == nil {
		t.Error("expected error while getting non-existent event, got none")
	}
}
