package main

import (
	"os"
	"os/exec"
)

func ExecuteCommand(env []EnvVar, command []string) error {
	for _, envVar := range env {
		_, ok := os.LookupEnv(envVar.Name)
		if ok {
			err := os.Unsetenv(envVar.Name)
			if err != nil {
				return err
			}
		}

		if envVar.MustBeRemoved {
			continue
		}

		err := os.Setenv(envVar.Name, envVar.Value)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command(command[0], command[1:]...) //nolint: gosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
