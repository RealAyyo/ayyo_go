package config

import (
	"errors"
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.

var (
	ErrConfigNotFound   = errors.New("config not found")
	ErrFailedReadConfig = errors.New("failed to read config")
)

type Config struct {
	Logger    LoggerConf    `yaml:"logger"`
	Storage   Storage       `yaml:"storage"`
	DB        DBConf        `yaml:"db"`
	HTTP      HTTPConf      `yaml:"http"`
	GRPC      GRPCConf      `yaml:"grpc"`
	RabbitMQ  RabbitMQConf  `yaml:"rabbit"`
	Scheduler SchedulerConf `yaml:"scheduler"`
	Env       string        `yaml:"env"  env-default:"local"`
}

type SchedulerConf struct {
	Interval string `yaml:"interval" env-default:"5m"`
}

type DBConf struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	Database string `yaml:"database" env-required:"true"`
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
}

type HTTPConf struct {
	Host string `yaml:"host" env-default:"0.0.0.0"`
	Port string `yaml:"port" env-default:"8888"`
}

type RabbitMQConf struct {
	User     string `yaml:"user" env-required:"true"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5672"`
	Password string `yaml:"password" env-required:"true"`
}

type GRPCConf struct {
	Port string `yaml:"port" env-default:"50051"`
}

type Storage struct {
	Type string `yaml:"type" env-default:"MEMORY"`
}

type LoggerConf struct {
	Level string `yaml:"level" env-default:"INFO"`
}

func NewConfig() (*Config, error) {
	configPath := fetchConfigPath()
	if configPath == "" {
		return nil, ErrConfigNotFound
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, ErrConfigNotFound
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, ErrFailedReadConfig
	}
	return &cfg, nil
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
