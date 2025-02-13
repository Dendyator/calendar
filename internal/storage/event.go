package storage

import (
	"time"

	"github.com/google/uuid" //nolint
)

type Event struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	StartTime   time.Time `db:"start_time"`
	EndTime     time.Time `db:"end_time"`
	UserID      uuid.UUID `db:"user_id"`
}

type Interface interface {
	CreateEvent(event Event) error
	UpdateEvent(id uuid.UUID, newEvent Event) error
	DeleteEvent(id uuid.UUID) error
	GetEvent(id uuid.UUID) (Event, error)
	ListEvents() ([]Event, error)
	ListEventsByDay(date time.Time) ([]Event, error)
	ListEventsByWeek(start time.Time) ([]Event, error)
	ListEventsByMonth(start time.Time) ([]Event, error)
	DeleteOldEvents(before time.Time) error
}
