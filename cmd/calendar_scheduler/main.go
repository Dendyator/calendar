package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/Dendyator/calendar/internal/config"                 //nolint
	"github.com/Dendyator/calendar/internal/logger"                 //nolint
	"github.com/Dendyator/calendar/internal/rabbitmq"               //nolint
	sqlstorage "github.com/Dendyator/calendar/internal/storage/sql" //nolint
	"github.com/google/uuid"                                        //nolint
	_ "github.com/lib/pq"                                           //nolint
)

type Notification struct {
	EventID   uuid.UUID `json:"eventId"`
	Title     string    `json:"title"`
	StartTime int64     `json:"startTime"`
}

func main() {
	configPath := flag.String("config", "configs/scheduler_config.yaml",
		"Path to configuration file")
	flag.Parse()

	cfg := config.LoadConfig(*configPath)
	logg := logger.New(cfg.Logger.Level)

	logg.Info("Starting scheduler...")

	rabbit, err := rabbitmq.New(cfg.RabbitMQ.DSN, logg)
	if err != nil {
		logg.Error("Failed to connect to RabbitMQ: " + err.Error())
		return
	}
	defer rabbit.Close()

	logg.Info("Connected to RabbitMQ")

	err = rabbit.DeclareQueue("notifications")
	if err != nil {
		logg.Error("Failed to declare RabbitMQ queue: " + err.Error())
		return
	}
	logg.Info("Declared RabbitMQ queue: notifications")

	store, err := sqlstorage.New(cfg.Database.DSN)
	if err != nil {
		logg.Error("Failed to connect to database: " + err.Error())
		return
	}
	logg.Info("Connected to database")

	for {
		events, err := store.ListEvents()
		if err != nil {
			logg.Error("Failed to list events: " + err.Error())
			continue
		}

		logg.Info(fmt.Sprintf("Retrieved %d events from database", len(events)))

		for _, event := range events {
			logg.Info(fmt.Sprintf("Processing event: %s at %s", event.Title, event.StartTime))

			timeUntilStart := time.Until(event.StartTime)
			logg.Info(fmt.Sprintf("Time until event starts: %v", timeUntilStart))

			if timeUntilStart < 24*time.Hour {
				notification := Notification{
					EventID:   event.ID,
					Title:     event.Title,
					StartTime: event.StartTime.Unix(),
				}

				body, err := json.Marshal(notification)
				if err != nil {
					logg.Error("Failed to marshal notification: " + err.Error())
					continue
				}

				err = rabbit.Publish("notifications", body)
				if err != nil {
					logg.Error("Failed to publish notification: " + err.Error())
				} else {
					logg.Info("Successfully published notification for event: " + notification.Title)
				}
			} else {
				logg.Info("Event does not require notification at this time")
			}
		}

		err = store.DeleteOldEvents(time.Now().AddDate(-1, 0, 0))
		if err != nil {
			logg.Error("Failed to delete old events: " + err.Error())
		} else {
			logg.Info("Old events deleted successfully")
		}

		logg.Info(fmt.Sprintf("Sleeping for %v", cfg.Scheduler.Interval))
		time.Sleep(cfg.Scheduler.Interval)
	}
}
