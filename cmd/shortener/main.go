package main

import (
	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/logging"
	"github.com/maxsnegir/url-shortener/internal/server"
	"github.com/maxsnegir/url-shortener/internal/storages"
)

func main() {
	cfg := config.NewConfig("configs/config.yaml")
	logger := logging.NewLogger(cfg.Logger.LogLevel)
	db := storages.NewURLDateBase()
	s := server.NewServer(cfg, logger, db)
	if err := s.Start(); err != nil {
		s.Logger.Error(err)
	}
}
