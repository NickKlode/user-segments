package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	HTTP `yaml:"http"`
	Log  `yaml:"log"`
	Postgres
}

type HTTP struct {
	Port    string        `yaml:"port" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
}

type Log struct {
	Level string `yaml:"level" env-default:"debug"`
}

type Postgres struct {
	URL string `env-required:"true" env:"POSTGRES_URL"`
}

func NewConfig() *Config {

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading .env file. %v", err)
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		logrus.Fatalf("config path is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logrus.Fatalf("config file does not exist. %v", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		logrus.Fatalf("cannot read config. %v", err)
	}

	return &cfg
}
