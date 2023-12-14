package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestEnvDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "testenv")
	require.NoError(t, err)
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0666); err != nil {
			require.NoError(t, err)
		}
	}
	return dir
}

func TestReadDir(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected []EnvVar
		wantErr  bool
	}{
		{
			name: "Normal variables",
			files: map[string]string{
				"VAR1": "value1",
				"VAR2": "value2",
			},
			expected: []EnvVar{
				{Name: "VAR1", Value: "value1"},
				{Name: "VAR2", Value: "value2"},
			},
			wantErr: false,
		},
		{
			name: "Empty variable",
			files: map[string]string{
				"EMPTY_VAR": "",
			},
			expected: []EnvVar{
				{Name: "EMPTY_VAR", Value: ""},
			},
			wantErr: false,
		},
		{
			name: "Variable with newlines",
			files: map[string]string{
				"NEWLINE_VAR": "line1\x00line2",
			},
			expected: []EnvVar{
				{Name: "NEWLINE_VAR", Value: "line1\nline2"},
			},
			wantErr: false,
		},
		{
			name:     "No files",
			files:    map[string]string{},
			expected: nil,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := createTestEnvDir(t, tt.files)
			defer os.RemoveAll(dir) // чистим за собой

			env, err := ReadEnv(dir)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, env)
			}
		})
	}
}
