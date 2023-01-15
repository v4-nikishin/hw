package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/storage/memory"
)

func TestAPI(t *testing.T) {
	event := storage.Event{
		UUID:  "123e4567-e89b-12d3-a456-426655440000",
		Title: "Event title",
		User:  "Event user",
		Date:  "2023-01-10",
		Begin: "20:00:00",
		End:   "21:30:00",
	}

	logg := logger.New(config.LoggerConf{Level: "debug"}, os.Stdout)
	repo := memorystorage.New()
	calendar := app.New(logg, repo)
	server := NewServer(config.ServerHTTP{Host: "localhost", Port: "8080"}, logg, calendar)

	ctx := context.Background()
	go func() {
		server.Start(ctx)
	}()
	defer server.Stop(ctx)

	cases := []struct {
		name         string
		method       string
		target       string
		event        *storage.Event
		responseCode int
	}{
		{"not found", http.MethodGet, "/qqq", nil, http.StatusNotFound},
		{"create event", http.MethodPost, "/create", &event, http.StatusOK},
		{"get event", http.MethodGet, "/event", &event, http.StatusOK},
		{"get invalid event", http.MethodGet, "/event", nil, http.StatusInternalServerError},
		{"update event", http.MethodPut, "/update", &event, http.StatusOK},
		{"invalid method", http.MethodPut, "/list", nil, http.StatusMethodNotAllowed},
		{"list events", http.MethodGet, "/list", nil, http.StatusOK},
		{"get events on date", http.MethodGet, "/events_on_date", &event, http.StatusOK},
		{"delete event", http.MethodDelete, "/delete", &event, http.StatusOK},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var output []byte
			if c.event != nil {
				o, err := json.Marshal(c.event)
				require.NoError(t, err)
				output = o
			}

			req, err := http.NewRequestWithContext(ctx, c.method, "http://localhost:8080"+c.target, bytes.NewBuffer(output))
			require.NoError(t, err)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer res.Body.Close()
			require.Equal(t, c.responseCode, res.StatusCode)
		})
	}
}
