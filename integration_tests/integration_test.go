package integration_tests_test

import (
	"encoding/json"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/logger"                 //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/rabbitmq"               //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage"                //nolint
	sqlstorage "github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage/sql" //nolint
	"github.com/google/uuid"                                                        //nolint
	_ "github.com/lib/pq"                                                           //nolint
	"github.com/stretchr/testify/assert"                                            //nolint
)

func TestCreateAndDuplicateEvent(t *testing.T) {
	stor, err := sqlstorage.New("user=user password=password dbname=calendar host=db port=5432 sslmode=disable")
	assert.NoError(t, err)
	assert.NotNil(t, stor)

	event := storage.Event{
		ID:          uuid.New(),
		Title:       "Test Event",
		Description: "Description of the test event",
		StartTime:   time.Now().Add(1 * time.Hour),
		EndTime:     time.Now().Add(2 * time.Hour),
		UserID:      uuid.New(),
	}

	err = stor.CreateEvent(event)
	assert.NoError(t, err)
	log.Println("Event created:", event)

	err = stor.CreateEvent(event)
	if assert.Error(t, err) {
		assert.True(t, strings.Contains(err.Error(), "duplicate key value"), "Error message should indicate duplicate key")
	}
	log.Println("Attempting to create duplicate event resulted in expected error:", err.Error())
}

func TestListEvents(t *testing.T) {
	stor, err := sqlstorage.New("user=user password=password dbname=calendar host=postgres_db port=5432 sslmode=disable")
	assert.NoError(t, err)
	assert.NotNil(t, stor)

	now := time.Now()
	nextWeek := now.Add(7 * 24 * time.Hour)
	twoWeeksLater := now.Add(14 * 24 * time.Hour)

	eventToday := storage.Event{
		ID:          uuid.New(),
		Title:       "Event Today",
		Description: "This event is scheduled for today.",
		StartTime:   now.Add(1 * time.Hour),
		EndTime:     now.Add(2 * time.Hour),
		UserID:      uuid.New(),
	}
	assert.NoError(t, stor.CreateEvent(eventToday))

	eventThisWeek := storage.Event{
		ID:          uuid.New(),
		Title:       "Event This Week",
		Description: "This event is scheduled within this week.",
		StartTime:   nextWeek.Add(-2 * 24 * time.Hour), // через 5 дней
		EndTime:     nextWeek.Add(-2*24*time.Hour + 1*time.Hour),
		UserID:      uuid.New(),
	}
	assert.NoError(t, stor.CreateEvent(eventThisWeek))

	eventNextWeek := storage.Event{
		ID:          uuid.New(),
		Title:       "Event Next Week",
		Description: "This event is scheduled for the next two weeks.",
		StartTime:   twoWeeksLater.Add(1 * time.Hour),
		EndTime:     twoWeeksLater.Add(2 * time.Hour),
		UserID:      uuid.New(),
	}
	assert.NoError(t, stor.CreateEvent(eventNextWeek))

	events, err := stor.ListEventsByDay(now)
	assert.NoError(t, err)
	log.Printf("Events for the day: %v\n", events)

	events, err = stor.ListEventsByWeek(now)
	assert.NoError(t, err)
	log.Printf("Events for the week: %v\n", events)

	events, err = stor.ListEventsByMonth(now)
	assert.NoError(t, err)
	log.Printf("Events for the month: %v\n", events)
}

type Notification struct {
	EventID   uuid.UUID `json:"eventId"`
	Title     string    `json:"title"`
	StartTime int64     `json:"startTime"`
}

type NotificationStatus struct {
	EventID uuid.UUID `json:"eventId"`
	Status  string    `json:"status"`
	Details string    `json:"details"`
}

func TestProcessNotificationIntegration(t *testing.T) {
	logg := logger.New("info")

	rabbit, err := rabbitmq.New("amqp://guest:guest@rabbitmq:5672/", logg)
	assert.NoError(t, err)
	defer rabbit.Close()

	err = rabbit.DeclareQueue("notifications")
	assert.NoError(t, err)

	err = rabbit.DeclareQueue("notification_statuses")
	assert.NoError(t, err)

	notification := Notification{
		EventID:   uuid.New(),
		Title:     "Integration Test Event",
		StartTime: time.Now().Unix(),
	}
	body, err := json.Marshal(notification)
	assert.NoError(t, err)

	err = rabbit.Publish("notifications", body)
	assert.NoError(t, err)

	deliveries, err := rabbit.Consume("notification_statuses")
	assert.NoError(t, err)

	timeout := time.After(3 * time.Second)
	done := make(chan bool)
	go func() {
		for msg := range deliveries {
			var status NotificationStatus
			err := json.Unmarshal(msg.Body, &status)
			assert.NoError(t, err)

			if status.EventID == notification.EventID {
				assert.Equal(t, "processed", status.Status)
				assert.Equal(t, "Notification processed successfully", status.Details)
				log.Printf("Notification processed successfully for event ID: %s", status.EventID)
				done <- true
				break
			}
		}
	}()

	select {
	case <-done:
	case <-timeout:
		t.Fatal("Test timed out waiting for notification status")
	}
}
