package sqlstorage

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/storage"
)

func TestStorage(t *testing.T) {
	// t.Skip()
	ctx := context.Background()
	logg := logger.New(config.LoggerConf{Level: logger.DebugStr}, os.Stdout)
	s, err := New(ctx,
		config.SQLConf{DSN: "host=localhost port=5432 user=postgres password=postgres dbname=calendar sslmode=disable"},
		logg)
	require.NoError(t, err)
	defer s.Close()

	t.Run("check insert", func(t *testing.T) {
		err := s.CreateEvent(storage.Event{UUID: "UUID", Title: "TITLE"})
		require.NoError(t, err)
	})
	t.Run("check list", func(t *testing.T) {
		events, err := s.Events()
		require.Equal(t, len(events), 1)
		require.NoError(t, err)
	})
	t.Run("check update event", func(t *testing.T) {
		err := s.UpdateEvent("UUID", "TITLE1")
		require.NoError(t, err)
	})
	t.Run("check get event", func(t *testing.T) {
		e, err := s.GetEvent("UUID")
		require.NoError(t, err)
		require.Equal(t, e.UUID, "UUID")
		require.Equal(t, e.Title, "TITLE1")
	})
	t.Run("check delete", func(t *testing.T) {
		err := s.DeleteEvent("UUID")
		require.NoError(t, err)
	})
}
