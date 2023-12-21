package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadEnv(t *testing.T) {
	t.Run("reads environment variables from files", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "envdir")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		err = os.WriteFile(filepath.Join(dir, "VAR1"), []byte("value1\n"), 0o600)
		require.NoError(t, err)

		err = os.WriteFile(filepath.Join(dir, "VAR2"), []byte("value2\n"), 0o600)
		require.NoError(t, err)

		envVars, err := ReadEnv(dir)
		require.NoError(t, err)

		require.Len(t, envVars, 2)
		require.Contains(t, envVars, EnvVar{Name: "VAR1", Value: "value1"})
		require.Contains(t, envVars, EnvVar{Name: "VAR2", Value: "value2"})
	})

	t.Run("ignores directories and files with =", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "envdir")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		err = os.Mkdir(filepath.Join(dir, "DIR"), 0o755)
		require.NoError(t, err)

		err = os.WriteFile(filepath.Join(dir, "VAR=3"), []byte("value3\n"), 0o600)
		require.NoError(t, err)

		envVars, err := ReadEnv(dir)
		require.NoError(t, err)

		require.Empty(t, envVars)
	})

	t.Run("trims spaces and tabs from the end of values", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "envdir")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		err = os.WriteFile(filepath.Join(dir, "VAR4"), []byte("value4 \t\n"), 0o600)
		require.NoError(t, err)

		envVars, err := ReadEnv(dir)
		require.NoError(t, err)

		require.Len(t, envVars, 1)
		require.Contains(t, envVars, EnvVar{Name: "VAR4", Value: "value4"})
	})

	t.Run("replaces null bytes with newlines", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "envdir")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		err = os.WriteFile(filepath.Join(dir, "VAR5"), []byte("line1\x00line2\n"), 0o600)
		require.NoError(t, err)

		envVars, err := ReadEnv(dir)
		require.NoError(t, err)

		require.Len(t, envVars, 1)
		require.Contains(t, envVars, EnvVar{Name: "VAR5", Value: "line1\nline2"})
	})

	t.Run("sets MustBeRemoved for empty files", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "envdir")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		err = os.WriteFile(filepath.Join(dir, "VAR6"), []byte("\n"), 0o600)
		require.NoError(t, err)

		envVars, err := ReadEnv(dir)
		require.NoError(t, err)

		require.Len(t, envVars, 1)
		require.Contains(t, envVars, EnvVar{Name: "VAR6", Value: "", MustBeRemoved: true})
	})
}
