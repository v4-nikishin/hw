package app

import (
	"context"
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
	EventsOnDate(date string) ([]storage.Event, error)
	Close()
}

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

func New(logger *logger.Logger, storage Storage) *App {
	return &App{log: logger, storage: storage}
}

func (a *App) CreateEvent(e storage.Event) error {
	if a.isBusyDatetime(e) {
		return fmt.Errorf("datetime is busy")
	}
	return a.storage.CreateEvent(e)
}

func (a *App) GetEvent(id string) (storage.Event, error) {
	return a.storage.GetEvent(id)
}

func (a *App) UpdateEvent(id string, e storage.Event) error {
	if a.isBusyDatetime(e) {
		return fmt.Errorf("datetime is busy")
	}
	return a.storage.UpdateEvent(id, e)
}

func (a *App) DeleteEvent(id string) error {
	return a.storage.DeleteEvent(id)
}

func (a *App) Events() ([]storage.Event, error) {
	return a.storage.Events()
}

func (a *App) EventsOnDate(date string) ([]storage.Event, error) {
	return a.storage.EventsOnDate(date)
}

func (a *App) isBusyDatetime(e storage.Event) bool {
	var err error
	defer func() {
		if err != nil {
			a.log.Error("faled to handle event: " + err.Error())
		}
	}()
	events, err := a.EventsOnDate(e.Date)
	if err != nil {
		return false
	}
	const format = "2006-01-02 15:04:05"
	newBegin, err := time.Parse(format, e.Date+" "+e.Begin)
	if err != nil {
		return false
	}
	newEnd, err := time.Parse(format, e.Date+" "+e.End)
	if err != nil {
		return false
	}
	for _, evt := range events {
		if evt.Date != e.Date || evt.UUID == e.UUID {
			continue
		}
		begin, err := time.Parse(format, evt.Date+" "+evt.Begin)
		if err != nil {
			return false
		}
		end, err := time.Parse(format, evt.Date+" "+evt.End)
		if err != nil {
			return false
		}
		if ((newBegin.After(begin) || newBegin.Equal(begin)) && (newBegin.Before(end) || newBegin.Equal(end))) ||
			((newEnd.After(begin) || newEnd.Equal(begin)) && (newEnd.Before(end) || newEnd.Equal(end))) {
			return true
		}
	}
	return false
}
