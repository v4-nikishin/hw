package memorystorage

import (
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

func (s *Storage) CreateEvent(e storage.Event) {
	s.mu.Lock()
	s.events[e.ID] = &e
	s.mu.Unlock()
}

func (s *Storage) GetEvent(id string) (storage.Event, bool) {
	s.mu.Lock()
	e, ok := s.events[id]
	s.mu.Unlock()
	return *e, ok
}

func (s *Storage) UpdateEvent(id string, title string) bool {
	s.mu.Lock()
	e, ok := s.events[id]
	if ok {
		e.Title = title
	}
	s.mu.Unlock()
	return ok
}

func (s *Storage) DeleteEvent(id string) {
	s.mu.Lock()
	delete(s.events, id)
	s.mu.Unlock()
}

func (s *Storage) Events() []storage.Event {
	events := []storage.Event{}
	s.mu.Lock()
	for _, e := range s.events {
		events = append(events, *e)
	}
	s.mu.Unlock()
	return events
}
