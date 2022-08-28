package main

import (
	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/logging"
	"github.com/maxsnegir/url-shortener/internal/server"
)

func main() {
	cfg := config.NewConfig("configs/config.yaml")
	logger := logging.NewLogger(cfg.Logger.LogLevel)
	s := server.NewServer(cfg, logger)
	if err := s.Start(); err != nil {
		s.Logger.Error(err)
	}
}
