package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("invalid path", func(t *testing.T) {
		_, err := ReadDir("qqq")
		require.NotNil(t, err)
	})
	t.Run("valid case", func(t *testing.T) {
		wanted := make(Environment, 4)
		wanted["BAR"] = EnvValue{"bar", false}
		wanted["EMPTY"] = EnvValue{"", false}
		wanted["FOO"] = EnvValue{"   foo\nwith new line", false}
		wanted["HELLO"] = EnvValue{"\"hello\"", false}
		res, err := ReadDir("./testdata/env")
		require.Nil(t, err)
		require.Equal(t, res, wanted)
	})
}
