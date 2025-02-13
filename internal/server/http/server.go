package internalhttp

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/Dendyator/calendar/internal/logger"  //nolint
	"github.com/Dendyator/calendar/internal/storage" //nolint
	"github.com/google/uuid"                         //nolint
	"github.com/gorilla/mux"                         //nolint
)

type Server struct {
	httpServer *http.Server
}

type ServerConfig struct {
	Host string
	Port string
}

func NewServer(cfg ServerConfig, logg *logger.Logger, store storage.Interface) *Server {
	router := mux.NewRouter()

	logg.Info("Setting up routes...")
	router.HandleFunc("/events", listEventsHandler(store, logg)).Methods(http.MethodGet)
	router.HandleFunc("/events", createEventHandler(store, logg)).Methods(http.MethodPost)
	router.HandleFunc("/events/{id:[0-9]+}", getEventHandler(store, logg)).Methods(http.MethodGet)
	router.HandleFunc("/events/{id:[0-9]+}", updateEventHandler(store, logg)).Methods(http.MethodPut)
	router.HandleFunc("/events/{id:[0-9]+}", deleteEventHandler(store, logg)).Methods(http.MethodDelete)
	logg.Info("Routes set up completed!")

	srv := &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:           loggingMiddleware(logg)(router),
		ReadHeaderTimeout: 5 * time.Second,
	}
	return &Server{
		httpServer: srv,
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		_ = s.Stop(context.Background())
	}()
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func listEventsHandler(store storage.Interface, logg *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		logg.Info("Handling GET request for listing events")

		events, err := store.ListEvents()
		if err != nil {
			logg.Errorf("Failed to list events: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if len(events) == 0 {
			logg.Info("No events found")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("[]"))
			return
		}

		err = json.NewEncoder(w).Encode(events)
		if err != nil {
			logg.Errorf("Failed to encode events to JSON: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		logg.Info("Events successfully listed.")
	}
}

func createEventHandler(store storage.Interface, logg *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logg.Infof("Handling POST request")
		var event storage.Event
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			logg.Errorf("Failed to decode event: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if err := store.CreateEvent(event); err != nil {
			logg.Errorf("Failed to create event: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		logg.Infof("Event created: %s", event.ID)
		w.WriteHeader(http.StatusCreated)
	}
}

func getEventHandler(store storage.Interface, logg *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logg.Infof("Handling GET request for a single event")
		id := r.URL.Path[len("/events/"):]
		parse, err := uuid.Parse(id)
		if err != nil {
			return
		}
		event, err := store.GetEvent(parse)
		if err != nil {
			logg.Errorf("Failed to get event: %v", err)
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(event)
	}
}

func updateEventHandler(store storage.Interface, logg *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logg.Infof("Handling PUT request")
		id := r.URL.Path[len("/events/"):]
		parse, err := uuid.Parse(id)
		if err != nil {
			return
		}
		var event storage.Event
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			logg.Errorf("Failed to decode event: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if err := store.UpdateEvent(parse, event); err != nil {
			logg.Errorf("Failed to update event: %v", err)
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		logg.Infof("Event updated: %s", id)
		w.WriteHeader(http.StatusOK)
	}
}

func deleteEventHandler(store storage.Interface, logg *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logg.Infof("Handling DELETE request")
		id := r.URL.Path[len("/events/"):]
		parse, err := uuid.Parse(id)
		if err != nil {
			return
		}
		if err := store.DeleteEvent(parse); err != nil {
			logg.Errorf("Failed to delete event: %v", err)
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		logg.Infof("Event deleted: %s", id)
		w.WriteHeader(http.StatusOK)
	}
}
