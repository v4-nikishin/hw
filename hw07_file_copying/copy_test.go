package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("empty source and/or destination", func(t *testing.T) {
		err := Copy("", "", 0, 0)
		require.NotNil(t, err)

		err = Copy("", "qqq", 0, 0)
		require.NotNil(t, err)

		err = Copy("qqq", "", 0, 0)
		require.NotNil(t, err)
	})
	t.Run("offset exceeds file size", func(t *testing.T) {
		fileName := "testdata/out_offset0_limit0.txt"

		in, err := os.Open(fileName)
		require.Nil(t, err)
		defer in.Close()

		sfi, err := os.Stat(fileName)
		require.Nil(t, err)

		err = Copy(fileName, "out.txt", sfi.Size()+1, 0)
		require.NotNil(t, err)
	})
	t.Run("unsupported file", func(t *testing.T) {
		err := Copy(os.TempDir(), "out.txt", 0, 0)
		require.NotNil(t, err)
	})
	t.Run("the same file", func(t *testing.T) {
		fileName := "testdata/out_offset0_limit0.txt"
		err := Copy(fileName, fileName, 0, 0)
		require.NotNil(t, err)
	})
}
