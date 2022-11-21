package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("empty cmd", func(t *testing.T) {
		ret := RunCmd(nil, nil)
		require.Equal(t, ret, 1)
	})
	t.Run("invalid cmd", func(t *testing.T) {
		ret := RunCmd([]string{"qqq"}, nil)
		require.Equal(t, ret, 1)
	})
}
