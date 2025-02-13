package memorystorage

import (
	"testing"
	"time"

	"github.com/Dendyator/calendar/internal/storage" //nolint:depguard
	"github.com/google/uuid"                         //nolint
	"github.com/stretchr/testify/assert"             //nolint
)

func TestStorage_CreateEvent(t *testing.T) {
	s := New()
	event := storage.Event{
		ID:          uuid.New(),
		Title:       "Test Event",
		Description: "This is a test event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	err := s.CreateEvent(event)
	assert.NoError(t, err)

	storedEvent, err := s.GetEvent(event.ID)
	assert.NoError(t, err)
	assert.Equal(t, event, storedEvent)
}

func TestStorage_CreateEvent_AlreadyExists(t *testing.T) {
	s := New()
	event := storage.Event{
		ID:          uuid.New(),
		Title:       "Test Event",
		Description: "This is a test event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	err := s.CreateEvent(event)
	assert.NoError(t, err)

	err = s.CreateEvent(event)
	assert.Error(t, err)
	assert.Equal(t, "event already exists", err.Error())
}

func TestStorage_UpdateEvent(t *testing.T) {
	s := New()
	event := storage.Event{
		ID:          uuid.New(),
		Title:       "Test Event",
		Description: "This is a test event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	err := s.CreateEvent(event)
	assert.NoError(t, err)

	updatedEvent := event
	updatedEvent.Title = "Updated Event"
	err = s.UpdateEvent(event.ID, updatedEvent)
	assert.NoError(t, err)

	storedEvent, err := s.GetEvent(event.ID)
	assert.NoError(t, err)
	assert.Equal(t, updatedEvent.Title, storedEvent.Title)
}

func TestStorage_UpdateEvent_NotFound(t *testing.T) {
	s := New()
	event := storage.Event{
		ID:          uuid.New(),
		Title:       "Test Event",
		Description: "This is a test event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	err := s.UpdateEvent(event.ID, event)
	assert.Error(t, err)
	assert.Equal(t, "event not found", err.Error())
}

func TestStorage_DeleteEvent(t *testing.T) {
	s := New()
	event := storage.Event{
		ID:          uuid.New(),
		Title:       "Test Event",
		Description: "This is a test event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	err := s.CreateEvent(event)
	assert.NoError(t, err)

	err = s.DeleteEvent(event.ID)
	assert.NoError(t, err)

	_, err = s.GetEvent(event.ID)
	assert.Error(t, err)
	assert.Equal(t, "event not found", err.Error())
}

func TestStorage_DeleteEvent_NotFound(t *testing.T) {
	s := New()
	eventID := uuid.New()

	err := s.DeleteEvent(eventID)
	assert.Error(t, err)
	assert.Equal(t, "event not found", err.Error())
}

func TestStorage_ListEvents(t *testing.T) {
	s := New()
	event1 := storage.Event{
		ID:          uuid.New(),
		Title:       "Test Event 1",
		Description: "This is the first test event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	event2 := storage.Event{
		ID:          uuid.New(),
		Title:       "Test Event 2",
		Description: "This is the second test event",
		StartTime:   time.Now().Add(2 * time.Hour),
		EndTime:     time.Now().Add(3 * time.Hour),
		UserID:      uuid.New(),
	}

	err := s.CreateEvent(event1)
	assert.NoError(t, err)
	err = s.CreateEvent(event2)
	assert.NoError(t, err)

	events, err := s.ListEvents()
	assert.NoError(t, err)
	assert.Len(t, events, 2)
}

func TestStorage_ListEventsByDay(t *testing.T) {
	s := New()
	event := storage.Event{
		ID:          uuid.New(),
		Title:       "Daily Event",
		Description: "This event is today",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	err := s.CreateEvent(event)
	assert.NoError(t, err)

	events, err := s.ListEventsByDay(time.Now())
	assert.NoError(t, err)
	assert.Len(t, events, 1)
}

func TestStorage_DeleteOldEvents(t *testing.T) {
	s := New()
	oldEvent := storage.Event{
		ID:          uuid.New(),
		Title:       "Old Event",
		Description: "This event is old",
		StartTime:   time.Now().Add(-2 * time.Hour),
		EndTime:     time.Now().Add(-1 * time.Hour),
		UserID:      uuid.New(),
	}

	newEvent := storage.Event{
		ID:          uuid.New(),
		Title:       "New Event",
		Description: "This event is new",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.New(),
	}

	err := s.CreateEvent(oldEvent)
	assert.NoError(t, err)
	err = s.CreateEvent(newEvent)
	assert.NoError(t, err)

	err = s.DeleteOldEvents(time.Now())
	assert.NoError(t, err)

	_, err = s.GetEvent(oldEvent.ID)
	assert.Error(t, err)
	assert.Equal(t, "event not found", err.Error())

	storedNewEvent, err := s.GetEvent(newEvent.ID)
	assert.NoError(t, err)
	assert.Equal(t, newEvent, storedNewEvent)
}
