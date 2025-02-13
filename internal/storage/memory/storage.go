package memorystorage

import (
	"errors"
	"sync"
	"time"

	"github.com/Dendyator/calendar/internal/storage" //nolint:depguard
	"github.com/google/uuid"                         //nolint
)

type Storage struct {
	mu     sync.RWMutex
	events map[uuid.UUID]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[uuid.UUID]storage.Event),
	}
}

func (s *Storage) CreateEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; exists {
		return errors.New("event already exists")
	}
	s.events[event.ID] = event

	return nil
}

func (s *Storage) UpdateEvent(id uuid.UUID, newEvent storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return errors.New("event not found")
	}
	s.events[id] = newEvent

	return nil
}

func (s *Storage) DeleteEvent(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return errors.New("event not found")
	}
	delete(s.events, id)

	return nil
}

func (s *Storage) GetEvent(id uuid.UUID) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, exists := s.events[id]

	if !exists {
		return event, errors.New("event not found")
	}
	return event, nil
}

func (s *Storage) ListEvents() ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]storage.Event, 0, len(s.events))

	for _, event := range s.events {
		events = append(events, event)
	}
	return events, nil
}

func (s *Storage) DeleteOldEvents(before time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, event := range s.events {
		if event.EndTime.Before(before) {
			delete(s.events, id)
		}
	}

	return nil
}

func (s *Storage) ListEventsByDay(date time.Time) ([]storage.Event, error) {
	start := date.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)
	var events []storage.Event
	for _, event := range s.events {
		if event.StartTime.After(start) && event.StartTime.Before(end) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (s *Storage) ListEventsByWeek(start time.Time) ([]storage.Event, error) {
	end := start.AddDate(0, 0, 7)
	var events []storage.Event
	for _, event := range s.events {
		if event.StartTime.After(start) && event.StartTime.Before(end) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (s *Storage) ListEventsByMonth(start time.Time) ([]storage.Event, error) {
	end := start.AddDate(0, 1, 0)
	var events []storage.Event
	for _, event := range s.events {
		if event.StartTime.After(start) && event.StartTime.Before(end) {
			events = append(events, event)
		}
	}
	return events, nil
}
