package app

import (
	"context"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(message string)
	Error(message string)
}

type Storage interface {
	CreateEvent(ctx context.Context, id, title string) error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	a.logger.Info("Creating event: " + title)
	return a.storage.CreateEvent(ctx, id, title)
}
