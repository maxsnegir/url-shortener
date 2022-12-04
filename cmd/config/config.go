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
	DatabaseDsn     = ""
)

// Config общие настройки для сервиса
type Config struct {
	Server struct {
		ServerAddress string
	}
	Logger struct {
		LogLevel string
	}
	Shortener struct {
		BaseURL string
	}
	Authorization struct {
		SecretKey string `env:"SECRET_KEY" envDefault:"super_secret"`
	}
	Storage struct {
		FileStoragePath string
		DatabaseDSN     string
	}
}

// NewConfig создание нового конфига настроек
func NewConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	// Server
	flag.StringVar(&cfg.Server.ServerAddress, "a", utils.GetEnv("SERVER_ADDRESS", ServerAddress), "server address")
	// Logger
	flag.StringVar(&cfg.Logger.LogLevel, "l", utils.GetEnv("LOG_LEVEL", LogLevel), "set log level")
	// Shortener
	flag.StringVar(&cfg.Shortener.BaseURL, "b", utils.GetEnv("BASE_URL", BaseURL), "base shortener address")
	// Storage
	flag.StringVar(&cfg.Storage.FileStoragePath, "f", utils.GetEnv("FILE_STORAGE_PATH", FileStoragePath), "name of file storage")
	flag.StringVar(&cfg.Storage.DatabaseDSN, "d", utils.GetEnv("DATABASE_DSN", DatabaseDsn), "db dsn")
	flag.Parse()
	return cfg, nil
}
