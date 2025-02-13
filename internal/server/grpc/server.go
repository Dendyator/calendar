package grpc

import (
	"context"
	"time"

	pb "github.com/Dendyator/calendar/api/pb"        //nolint
	"github.com/Dendyator/calendar/internal/logger"  //nolint
	"github.com/Dendyator/calendar/internal/storage" //nolint
	"github.com/google/uuid"                         //nolint
	_ "github.com/jackc/pgx/v4/stdlib"               //nolint
)

type Server struct {
	pb.UnimplementedEventServiceServer
	storage storage.Interface
	logg    *logger.Logger
}

func NewGRPCServer(storage storage.Interface, logg *logger.Logger) *Server {
	return &Server{storage: storage, logg: logg}
}

func (s *Server) CreateEvent(_ context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	s.logg.Info("Creating event: " + req.GetEvent().Title)
	event := storage.Event{
		ID:          uuid.New(),
		Title:       req.GetEvent().Title,
		Description: req.GetEvent().Description,
		StartTime:   time.Unix(req.GetEvent().StartTime, 0),
		EndTime:     time.Unix(req.GetEvent().EndTime, 0),
		UserID:      uuid.MustParse(req.GetEvent().UserId),
	}
	err := s.storage.CreateEvent(event)
	if err != nil {
		s.logg.Error("Failed to create event: " + err.Error())
	}
	return &pb.CreateEventResponse{}, err
}

func (s *Server) UpdateEvent(_ context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	s.logg.Info("Updating event ID: " + req.GetId())
	newEvent := storage.Event{
		ID:          uuid.MustParse(req.GetEvent().Id),
		Title:       req.GetEvent().Title,
		Description: req.GetEvent().Description,
		StartTime:   time.Unix(req.GetEvent().StartTime, 0),
		EndTime:     time.Unix(req.GetEvent().EndTime, 0),
		UserID:      uuid.MustParse(req.GetEvent().UserId),
	}
	err := s.storage.UpdateEvent(newEvent.ID, newEvent)
	if err != nil {
		s.logg.Error("Failed to update event: " + err.Error())
	}
	return &pb.UpdateEventResponse{}, err
}

func (s *Server) DeleteEvent(_ context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	s.logg.Info("Deleting event ID: " + req.GetId())
	err := s.storage.DeleteEvent(uuid.MustParse(req.GetId()))
	if err != nil {
		s.logg.Error("Failed to delete event: " + err.Error())
	}
	return &pb.DeleteEventResponse{}, err
}

func (s *Server) GetEvent(_ context.Context, req *pb.GetEventRequest) (*pb.GetEventResponse, error) {
	s.logg.Info("Retrieving event ID: " + req.GetId())
	event, err := s.storage.GetEvent(uuid.MustParse(req.GetId()))
	if err != nil {
		s.logg.Error("Failed to get event: " + err.Error())
		return nil, err
	}
	return &pb.GetEventResponse{
		Event: &pb.Event{
			Id:          event.ID.String(),
			Title:       event.Title,
			Description: event.Description,
			StartTime:   event.StartTime.Unix(),
			EndTime:     event.EndTime.Unix(),
			UserId:      event.UserID.String(),
		},
	}, nil
}

func (s *Server) ListEvents(_ context.Context, _ *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	s.logg.Info("Listing all events")
	events, err := s.storage.ListEvents()
	if err != nil {
		s.logg.Error("Failed to list events: " + err.Error())
		return nil, err
	}
	pbEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		pbEvents[i] = &pb.Event{
			Id:          event.ID.String(),
			Title:       event.Title,
			Description: event.Description,
			StartTime:   event.StartTime.Unix(),
			EndTime:     event.EndTime.Unix(),
			UserId:      event.UserID.String(),
		}
	}
	return &pb.ListEventsResponse{Events: pbEvents}, nil
}

func (s *Server) ListEventsByDay(_ context.Context, req *pb.ListEventsByDayRequest,
) (*pb.ListEventsByDayResponse, error) {
	date := time.Unix(req.GetDate(), 0)
	events, err := s.storage.ListEventsByDay(date)
	if err != nil {
		s.logg.Error("Failed to list events by day: " + err.Error())
		return nil, err
	}
	return &pb.ListEventsByDayResponse{Events: convertToPBEvents(events)}, nil
}

func (s *Server) ListEventsByWeek(_ context.Context, req *pb.ListEventsByWeekRequest,
) (*pb.ListEventsByWeekResponse, error) {
	start := time.Unix(req.GetStart(), 0)
	events, err := s.storage.ListEventsByWeek(start)
	if err != nil {
		s.logg.Error("Failed to list events by week: " + err.Error())
		return nil, err
	}
	return &pb.ListEventsByWeekResponse{Events: convertToPBEvents(events)}, nil
}

func (s *Server) ListEventsByMonth(_ context.Context, req *pb.ListEventsByMonthRequest,
) (*pb.ListEventsByMonthResponse, error) {
	start := time.Unix(req.GetStart(), 0)
	events, err := s.storage.ListEventsByMonth(start)
	if err != nil {
		s.logg.Error("Failed to list events by month: " + err.Error())
		return nil, err
	}
	return &pb.ListEventsByMonthResponse{Events: convertToPBEvents(events)}, nil
}

func convertToPBEvents(events []storage.Event) []*pb.Event {
	pbEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		pbEvents[i] = &pb.Event{
			Id:          event.ID.String(),
			Title:       event.Title,
			Description: event.Description,
			StartTime:   event.StartTime.Unix(),
			EndTime:     event.EndTime.Unix(),
			UserId:      event.UserID.String(),
		}
	}
	return pbEvents
}
