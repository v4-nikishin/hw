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
	GetEvent(id string) (storage.Event, error)
	UpdateEvent(id string, title string) error
	DeleteEvent(id string) error
	Events() ([]storage.Event, error)
}

func New(logger *logger.Logger, storage Storage) *App {
	return &App{log: logger}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return a.storage.CreateEvent(storage.Event{UUID: id, Title: title})
}

func (a *App) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	return a.storage.GetEvent(id)
}

func (a *App) UpdateEvent(ctx context.Context, id, title string) error {
	return a.storage.UpdateEvent(id, title)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(id)
}

func (a *App) Events() ([]storage.Event, error) {
	return a.storage.Events()
}
