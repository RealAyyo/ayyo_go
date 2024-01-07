package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf `yaml:"logger"`
	Env    string     `yaml:"env"  env-default:"local"`
}

type LoggerConf struct {
	Level string `yaml:"level" env-default:"INFO"`
}

func NewConfig() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("Path config is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("Config not found")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("Failed to read config: " + err.Error())
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

// TODO
