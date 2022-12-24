package memorystorage

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/storage"
)

func TestStorage(t *testing.T) {
	t.Skip()
	s := New()
	t.Run("check create", func(t *testing.T) {
		s.CreateEvent(storage.Event{UUID: "1", Title: "1"})
		require.Equal(t, s.events["1"].Title, "1")
		require.Equal(t, s.events["1"].UUID, "1")
	})
	t.Run("check get", func(t *testing.T) {
		e, err := s.GetEvent("1")
		require.Equal(t, e.Title, "1")
		require.Equal(t, e.UUID, "1")
		require.NoError(t, err)
	})
	t.Run("check update", func(t *testing.T) {
		err := s.UpdateEvent("1", "2")
		require.Equal(t, s.events["1"].Title, "2")
		require.NoError(t, err)
	})
	t.Run("check list", func(t *testing.T) {
		events, err := s.Events()
		require.Equal(t, len(events), 1)
		require.NoError(t, err)
	})
	t.Run("check delete", func(t *testing.T) {
		s.DeleteEvent("1")
		require.Equal(t, len(s.events), 0)
	})
}
