package main

import (
	"os"
	"strings"
)

type EnvVar struct {
	Name  string
	Value string
}

func ReadEnv(dir string) ([]EnvVar, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make([]EnvVar, 0, len(files))
	for _, file := range files {
		content, err := os.ReadFile(dir + "/" + file.Name())
		if err != nil {
			return nil, err
		}

		value := strings.ReplaceAll(string(content), "\x00", "\n")
		value = strings.TrimRight(value, "\n")

		env = append(env, EnvVar{
			Name:  file.Name(),
			Value: value,
		})
	}

	return env, nil
}
