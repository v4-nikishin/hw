package app

import (
	"fmt"
	"time"

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
	UpdateEvent(id string, e storage.Event) error
	DeleteEvent(id string) error
	Events() ([]storage.Event, error)
	Close()
}

func New(logger *logger.Logger, storage Storage) *App {
	return &App{log: logger, storage: storage}
}

func (a *App) CreateEvent(e storage.Event) error {
	evt, err := a.GetEvent(e.UUID)
	if err != nil {
		return a.storage.CreateEvent(e)
	}
	format := "2006-01-02 15:04:00"
	begin, err := time.Parse(format, evt.Date+" "+evt.Begin)
	if err != nil {
		return err
	}
	end, err := time.Parse(format, evt.Date+" "+evt.End)
	if err != nil {
		return err
	}
	newBegin, err := time.Parse(format, e.Date+" "+e.Begin)
	if err != nil {
		return err
	}
	newEnd, err := time.Parse(format, e.Date+" "+e.End)
	if err != nil {
		return err
	}

	if ((newBegin.After(begin) || newBegin.Equal(begin)) && (newBegin.Before(end) || newBegin.Equal(end))) ||
		((newEnd.After(begin) || newEnd.Equal(begin)) && (newEnd.Before(end) || newEnd.Equal(end))) {
		return fmt.Errorf("datetime is busy")
	}
	return a.storage.CreateEvent(e)
}

func (a *App) GetEvent(id string) (storage.Event, error) {
	return a.storage.GetEvent(id)
}

func (a *App) UpdateEvent(id string, e storage.Event) error {
	return a.storage.UpdateEvent(id, e)
}

func (a *App) DeleteEvent(id string) error {
	return a.storage.DeleteEvent(id)
}

func (a *App) Events() ([]storage.Event, error) {
	return a.storage.Events()
}
