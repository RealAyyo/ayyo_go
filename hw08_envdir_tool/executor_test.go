package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name    string
		env     []EnvVar
		command []string
		wantErr bool
	}{
		{
			name: "Execute echo with no env vars",
			env:  nil,
			command: []string{
				"echo",
				"Hello, world!",
			},
			wantErr: false,
		},
		{
			name: "Execute echo with env vars",
			env: []EnvVar{
				{Name: "TEST_VAR", Value: "123"},
			},
			command: []string{
				"bash",
				"-c",
				"echo $TEST_VAR",
			},
			wantErr: false,
		},
		{
			name: "Execute non-existing command",
			env:  nil,
			command: []string{
				"non_existing_command",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ExecuteCommand(tt.env, tt.command)
			if tt.wantErr {
				require.Error(t, err, "ExecuteCommand should return an error")
			} else {
				require.NoError(t, err, "ExecuteCommand should not return an error")
			}
		})
	}
}
