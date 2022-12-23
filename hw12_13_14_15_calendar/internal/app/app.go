package app

import (
	"context"

	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	log     *logger.Logger
	storage Storage
}

type Storage interface {
	CreateEvent(e storage.Event) error
	GetEvent(id string) (storage.Event, bool)
	UpdateEvent(id string, title string) bool
	DeleteEvent(id string) error
	Events() []storage.Event
}

func New(logger *logger.Logger, storage Storage) *App {
	return &App{log: logger}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

func (a *App) GetEvent(ctx context.Context, id string) (storage.Event, bool) {
	return a.storage.GetEvent(id)
}

func (a *App) UpdateEvent(ctx context.Context, id, title string) bool {
	return a.storage.UpdateEvent(id, title)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(id)
}

func (a *App) Events() []storage.Event {
	return a.storage.Events()
}
