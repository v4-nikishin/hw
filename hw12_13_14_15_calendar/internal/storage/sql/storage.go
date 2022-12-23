package sqlstorage

import (
	"context"

	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct { // TODO
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) CreateEvent(e storage.Event) error {
	return nil
}

func (s *Storage) GetEvent(id string) (storage.Event, bool) {
	return storage.Event{}, true
}

func (s *Storage) UpdateEvent(id string, title string) bool {
	return true
}

func (s *Storage) DeleteEvent(id string) error {
	return nil
}

func (s *Storage) Events() []storage.Event {
	return nil
}
