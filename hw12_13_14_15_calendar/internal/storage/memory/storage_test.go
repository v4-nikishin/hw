package memorystorage

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/storage"
)

func TestStorage(t *testing.T) {
	s := New()
	t.Run("check create", func(t *testing.T) {
		s.CreateEvent(storage.Event{ID: "1", Title: "1"})
		require.Equal(t, s.events["1"].Title, "1")
		require.Equal(t, s.events["1"].ID, "1")
	})
	t.Run("check get", func(t *testing.T) {
		e, ok := s.GetEvent("1")
		require.Equal(t, e.Title, "1")
		require.Equal(t, e.ID, "1")
		require.Equal(t, ok, true)
	})
	t.Run("check update", func(t *testing.T) {
		ok := s.UpdateEvent("1", "2")
		require.Equal(t, s.events["1"].Title, "2")
		require.Equal(t, ok, true)
	})
	t.Run("check list", func(t *testing.T) {
		events := s.Events()
		require.Equal(t, len(events), 1)
	})
	t.Run("check delete", func(t *testing.T) {
		s.DeleteEvent("1")
		require.Equal(t, len(s.events), 0)
	})
}
