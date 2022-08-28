package main

import (
	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/db"
	"github.com/maxsnegir/url-shortener/internal/logging"
	"github.com/maxsnegir/url-shortener/internal/server"
)

func main() {
	cfg := config.NewConfig("configs/config.yaml")
	logger := logging.NewLogger(cfg.Logger.LogLevel)
	redisClient, err := db.NewRedis(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	s := server.NewServer(cfg, logger, redisClient)
	if err := s.Start(); err != nil {
		s.Logger.Error(err)
	}
}
