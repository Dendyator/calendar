package sqlstorage

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/Dendyator/calendar/internal/storage" //nolint
	"github.com/google/uuid"                         //nolint
	_ "github.com/jackc/pgx/v4/stdlib"               //nolint
	"github.com/jmoiron/sqlx"                        //nolint
)

type Storage struct {
	DB *sqlx.DB
}

func New(dsn string) (*Storage, error) {
	log.Println("Using DSN:", dsn)
	var db *sqlx.DB
	var err error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Open("postgres", dsn)

		if err == nil {
			if err = db.Ping(); err == nil {
				log.Println("Successfully connected to the database!")
				return &Storage{DB: db}, nil
			}
		}
		log.Printf("Failed to connect to database: %v. Retrying...\n", err)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}

func (s *Storage) CreateEvent(event storage.Event) error {
	query := "INSERT INTO events (id, title, description, start_time," +
		" end_time, user_id) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := s.DB.Exec(query, event.ID, event.Title, event.Description, event.StartTime, event.EndTime, event.UserID)
	return err
}

func (s *Storage) UpdateEvent(id uuid.UUID, newEvent storage.Event) error {
	query := "UPDATE events SET title = $1, description = $2, start_time = $3, end_time = $4," +
		" user_id = $5 WHERE id = $6"
	res, err := s.DB.Exec(query, newEvent.Title, newEvent.Description, newEvent.StartTime,
		newEvent.EndTime, newEvent.UserID, id)
	if err != nil {
		return err
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return errors.New("event not found")
	}
	return nil
}

func (s *Storage) DeleteEvent(id uuid.UUID) error {
	query := "DELETE FROM events WHERE id = $1"
	res, err := s.DB.Exec(query, id)
	if err != nil {
		return err
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return errors.New("event not found")
	}
	return nil
}

func (s *Storage) GetEvent(id uuid.UUID) (storage.Event, error) {
	var event storage.Event
	query := "SELECT id, title, description, start_time, end_time, user_id FROM events WHERE id = $1"
	err := s.DB.Get(&event, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return event, errors.New("event not found")
	}
	return event, err
}

func (s *Storage) ListEvents() ([]storage.Event, error) {
	var events []storage.Event
	query := "SELECT id, title, description, start_time, end_time, user_id FROM events"
	err := s.DB.Select(&events, query)
	return events, err
}

func (s *Storage) DeleteOldEvents(before time.Time) error {
	query := "DELETE FROM events WHERE end_time < $1"
	_, err := s.DB.Exec(query, before)
	return err
}

func (s *Storage) ListEventsByDay(date time.Time) ([]storage.Event, error) {
	start := date.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)
	query := `SELECT id, title, description, start_time, end_time, user_id FROM events 
              WHERE start_time >= $1 AND start_time < $2`
	var events []storage.Event
	err := s.DB.Select(&events, query, start, end)
	return events, err
}

func (s *Storage) ListEventsByWeek(start time.Time) ([]storage.Event, error) {
	end := start.AddDate(0, 0, 7)
	query := `SELECT id, title, description, start_time, end_time, user_id FROM events
              WHERE start_time >= $1 AND start_time < $2`
	var events []storage.Event
	err := s.DB.Select(&events, query, start, end)
	return events, err
}

func (s *Storage) ListEventsByMonth(start time.Time) ([]storage.Event, error) {
	end := start.AddDate(0, 1, 0)
	query := `SELECT id, title, description, start_time, end_time, user_id FROM events
              WHERE start_time >= $1 AND start_time < $2`
	var events []storage.Event
	err := s.DB.Select(&events, query, start, end)
	return events, err
}
