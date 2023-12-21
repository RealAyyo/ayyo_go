package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-envdir /path/to/env/dir command arg1 arg2...")
		os.Exit(1)
	}

	envDir := os.Args[1]
	command := os.Args[2:]

	env, err := ReadEnv(envDir)
	if err != nil {
		fmt.Println("Error reading environment:", err)
		os.Exit(1)
	}

	err = ExecuteCommand(env, command)
	if err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
