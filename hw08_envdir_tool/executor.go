package main

import (
	"os"
	"os/exec"
)

func ExecuteCommand(env []EnvVar, command []string) error {
	cmd := exec.Command(command[0], command[1:]...)

	cmd.Env = os.Environ()
	for _, e := range env {
		cmd.Env = append(cmd.Env, e.Name+"="+e.Value)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
