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
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	logger := logging.NewLogger(cfg.Logger.LogLevel)
	urlDataBase := storages.NewURLDataBase()
	shortener := services.NewShortener(urlDataBase, cfg.Shortener.BaseURL)
	urlHandler := server.NewURLHandler(shortener, logger)
	s := server.NewServer(cfg, logger, urlHandler)
	logger.Fatal(s.Start())
}
