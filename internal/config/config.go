package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath = "configs/config.yaml"

type Config struct {
	HTTPServer HTTPServer `yaml:"http_server"`
	Postgres   Postgres   `yaml:"postgres"`
	Cache      Cache      `yaml:"cache"`
	Kafka      Kafka      `yaml:"kafka"`
}

type HTTPServer struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Cache struct {
	StartupSize int `yaml:"startup_size"`
}

type Kafka struct {
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Topic string `yaml:"topic"`
	Group string `yaml:"group"`
}

func New() *Config {
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return &cfg
}
