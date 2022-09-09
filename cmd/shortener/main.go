package main

import (
	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/logging"
	"github.com/maxsnegir/url-shortener/internal/server"
	"github.com/maxsnegir/url-shortener/internal/services"
	"github.com/maxsnegir/url-shortener/internal/storages"
	"log"
)

func main() {
	cfg := config.NewConfig("configs/config.yaml")
	logger := logging.NewLogger(cfg.Logger.LogLevel)
	urlDataBase := storages.NewURLDataBase()
	shortener := services.NewShortener(urlDataBase, cfg.Server.FullAddress)
	urlHandler := server.NewURLHandler(shortener, logger)
	s := server.NewServer(cfg, logger, urlHandler)
	log.Fatal(s.Start())
}
