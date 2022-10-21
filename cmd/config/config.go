package config

import (
	"github.com/caarlos0/env/v6"
)

// Config common settings for web application
type Config struct {
	Server struct {
		ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8000"`
	}
	Logger struct {
		LogLevel string `env:"LOG_LEVEL" envDefault:"DEBUG"`
	}
	Shortener struct {
		BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8000"`
		FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
	}
}

// NewConfig creates a new Config
func NewConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
