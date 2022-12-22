package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("check equal levels", func(t *testing.T) {
		require.Equal(t, levelNum(ErrorStr), Error)
		require.Equal(t, levelNum(WarnStr), Warn)
		require.Equal(t, levelNum(InfoStr), Info)
		require.Equal(t, levelNum(DebugStr), Debug)
	})
	t.Run("check less levels", func(t *testing.T) {
		require.Less(t, levelNum(ErrorStr), Warn)
		require.Less(t, levelNum(WarnStr), Info)
		require.Less(t, levelNum(InfoStr), Debug)
	})
	t.Run("check output levels", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(WarnStr, out)
		log.Error("Error")
		require.Contains(t, out.String(), ErrorTag)
		log.Warn("Warn")
		require.Contains(t, out.String(), WarnTag)
		log.Info("Info")
		require.NotContains(t, out.String(), InfoTag)
		log.Debug("Debug")
		require.NotContains(t, out.String(), DebugTag)
	})
}
