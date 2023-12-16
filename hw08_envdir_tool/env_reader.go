package main

import (
	"os"
	"path"
	"strings"
)

type EnvVar struct {
	Name          string
	Value         string
	MustBeRemoved bool
}

func ReadEnv(dir string) ([]EnvVar, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envVars := make([]EnvVar, 0, len(files))
	for _, file := range files {
		if file.IsDir() || strings.Contains(file.Name(), "=") {
			continue
		}

		content, err := os.ReadFile(path.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		lines := strings.Split(string(content), "\n")
		value := lines[0]
		value = strings.TrimRight(value, "\n")
		value = strings.ReplaceAll(value, string([]byte{0x00}), "\n")
		value = strings.TrimRight(value, " \t")

		envVars = append(envVars, EnvVar{
			Name:          file.Name(),
			Value:         value,
			MustBeRemoved: len(value) == 0,
		})
	}

	return envVars, nil
}
