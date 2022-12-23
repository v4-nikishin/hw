package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	cfg config.ServerConf
	log *logger.Logger
}

type Application interface { // TODO
}

func NewServer(cfg config.ServerConf, logger *logger.Logger, app Application) *Server {
	return &Server{cfg: cfg, log: logger}
}

func (s *Server) Start(ctx context.Context) error {
	addr := s.cfg.Host + ":" + s.cfg.Port
	s.log.Info("server started on address: " + addr)
	server := &http.Server{
		Addr:         addr,
		Handler:      loggingMiddleware(s, s.log),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	server.ListenAndServe()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("...calendar is stopped")
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Calendar!")
	case "/hello":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello world!")
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
