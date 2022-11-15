package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/maxsnegir/url-shortener/internal/utils"
)

const (
	ServerAddress   = "localhost:8080"
	LogLevel        = "DEBUG"
	BaseURL         = "http://localhost:8080"
	FileStoragePath = ""
)

// Config common settings for web application
type Config struct {
	Server struct {
		ServerAddress string
	}
	Logger struct {
		LogLevel string
	}
	Shortener struct {
		BaseURL         string
		FileStoragePath string
	}
	Authorization struct {
		SecretKey string `env:"SECRET_KEY" envDefault:"super_secret"`
	}
}

// NewConfig creates a new Config
func NewConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	flag.StringVar(&cfg.Server.ServerAddress, "a", utils.GetEnv("SERVER_ADDRESS", ServerAddress), "server address")
	flag.StringVar(&cfg.Shortener.BaseURL, "b", utils.GetEnv("BASE_URL", BaseURL), "base shortener address")
	flag.StringVar(&cfg.Shortener.FileStoragePath, "f", utils.GetEnv("FILE_STORAGE_PATH", FileStoragePath), "name of file storage")
	flag.StringVar(&cfg.Logger.LogLevel, "l", utils.GetEnv("LOG_LEVEL", LogLevel), "set log level")
	flag.Parse()
	return cfg, nil
}
