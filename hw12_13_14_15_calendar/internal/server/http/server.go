package internalhttp

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	cfg    config.ServerHTTP
	log    *logger.Logger
	server *http.Server
	app    *app.App
}

func NewServer(cfg config.ServerHTTP, logger *logger.Logger, app *app.App) *Server {
	return &Server{cfg: cfg, log: logger, app: app}
}

func (s *Server) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.cfg.Host, s.cfg.Port)
	s.log.Info("http server is started on address: " + addr)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      loggingMiddleware(s, s.log),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	s.log.Info("...http server is stopped")
	return nil
}

func (s *Server) fetchEvent(m string, w *http.ResponseWriter, r *http.Request) (storage.Event, bool) {
	if !s.allowedMethod(m, w, r) {
		return storage.Event{}, false
	}
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		s.log.Error("invalid request body: " + err.Error())
		return storage.Event{}, false
	}
	var e storage.Event
	err = json.Unmarshal(b, &e)
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		s.log.Error("failed to unmarshal request body: " + err.Error())
		return storage.Event{}, false
	}
	return e, true
}

func (s *Server) allowedMethod(m string, w *http.ResponseWriter, r *http.Request) bool {
	if m != r.Method {
		(*w).WriteHeader(http.StatusMethodNotAllowed)
		s.log.Error("method not allowed: " + m)
		return false
	}
	return true
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/create":
		e, ok := s.fetchEvent(http.MethodPost, &w, r)
		if !ok {
			return
		}
		if err := s.app.CreateEvent(e); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.log.Error("failed to get events: " + err.Error())
		}
	case "/event":
		e, ok := s.fetchEvent(http.MethodGet, &w, r)
		if !ok {
			return
		}
		e, err := s.app.GetEvent(e.UUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.log.Error("failed to get events: " + err.Error())
			return
		}
		s.writeEvent(w, &e)
	case "/update":
		e, ok := s.fetchEvent(http.MethodPut, &w, r)
		if !ok {
			return
		}
		if err := s.app.UpdateEvent(e.UUID, e); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.log.Error("failed to get events: " + err.Error())
		}
	case "/list":
		if !s.allowedMethod(http.MethodGet, &w, r) {
			return
		}
		evts, err := s.app.Events()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.log.Error("failed to get events" + err.Error())
			return
		}
		s.writeEvents(w, evts)
	case "/events_on_date":
		if !s.allowedMethod(http.MethodGet, &w, r) {
			return
		}
		if _, ok := s.fetchEvent(http.MethodGet, &w, r); !ok {
			return
		}
		evts, err := s.app.Events()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.log.Error("failed to get events on date" + err.Error())
			return
		}
		s.writeEvents(w, evts)
	case "/delete":
		e, ok := s.fetchEvent(http.MethodDelete, &w, r)
		if !ok {
			return
		}
		if err := s.app.DeleteEvent(e.UUID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.log.Error("failed to get events" + err.Error())
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Server) writeEvent(w http.ResponseWriter, resp *storage.Event) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		s.log.Error("response marshal error: " + err.Error())
	}
	_, err = w.Write(resBuf)
	if err != nil {
		s.log.Error("response marshal error: " + err.Error())
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (s *Server) writeEvents(w http.ResponseWriter, resp []storage.Event) {
	events := storage.Events{}
	events.Evets = resp
	resBuf, err := json.Marshal(events)
	if err != nil {
		s.log.Error("response marshal error: " + err.Error())
	}
	_, err = w.Write(resBuf)
	if err != nil {
		s.log.Error("response marshal error: " + err.Error())
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}
