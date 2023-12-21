package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("sets environment variables for command", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "envdir")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		err = os.WriteFile(filepath.Join(dir, "VAR1"), []byte("value1\n"), 0o600)
		require.NoError(t, err)

		err = os.WriteFile(filepath.Join(dir, "VAR2"), []byte("value2\n"), 0o600)
		require.NoError(t, err)

		envVars, err := ReadEnv(dir)
		require.NoError(t, err)

		cmd := []string{"printenv", "VAR1", "VAR2"}
		err = ExecuteCommand(envVars, cmd)
		require.NoError(t, err)
	})

	t.Run("unsets environment variables for command", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "envdir")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		err = os.WriteFile(filepath.Join(dir, "VAR1"), []byte("\n"), 0o600)
		require.NoError(t, err)

		envVars, err := ReadEnv(dir)
		require.NoError(t, err)

		cmd := []string{"printenv", "VAR1"}
		err = ExecuteCommand(envVars, cmd)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "exit status 1"))
	})
}
