package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(h http.Handler, log *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Info(fmt.Sprintf("%s %s %s %s %v", r.Method, r.URL.Path, r.Proto, r.UserAgent(), time.Since(start)))
	})
}
