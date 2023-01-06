package memorystorage

import (
	"fmt"
	"sync"

	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[string]*storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{events: make(map[string]*storage.Event)}
}

func (s *Storage) CreateEvent(e storage.Event) error {
	s.mu.RLock()
	s.events[e.UUID] = &e
	s.mu.RUnlock()
	return nil
}

func (s *Storage) GetEvent(id string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.events[id]
	if !ok {
		return storage.Event{}, fmt.Errorf("invalide event %s", id)
	}
	return *e, nil
}

func (s *Storage) UpdateEvent(id string, evt storage.Event) error {
	s.mu.RLock()
	e, ok := s.events[id]
	if ok {
		e.Title = evt.Title
	}
	s.mu.RUnlock()
	return nil
}

func (s *Storage) DeleteEvent(id string) error {
	s.mu.RLock()
	delete(s.events, id)
	s.mu.RUnlock()
	return nil
}

func (s *Storage) Events() ([]storage.Event, error) {
	events := []storage.Event{}
	s.mu.RLock()
	for _, e := range s.events {
		events = append(events, *e)
	}
	s.mu.RUnlock()
	return events, nil
}

func (s *Storage) Close() {
}
