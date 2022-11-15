package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/auth"
	"github.com/maxsnegir/url-shortener/internal/handlers"
	"github.com/maxsnegir/url-shortener/internal/logging"
	"github.com/maxsnegir/url-shortener/internal/server"
	"github.com/maxsnegir/url-shortener/internal/services"
	"github.com/maxsnegir/url-shortener/internal/storage"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	logger := logging.NewLogger(cfg.Logger.LogLevel)
	stor, err := storage.GetURLStorage(cfg)
	if err != nil {
		switch err.(type) {
		// Если случилась ошибка при загрузке данных из файла - все равно продолжим работу
		case storage.LoadingDumbDataError:
			logger.Error(err)
		default:
			logger.Fatal(err)
		}
	}
	shortener := services.NewShortener(stor, cfg.Shortener.BaseURL)
	authorization, err := auth.NewCookieAuthentication(cfg.Authorization.SecretKey)
	if err != nil {
		logger.Fatal(err)
	}
	urlHandler := handlers.NewURLHandler(shortener, authorization, logger)
	s := server.NewServer(cfg, logger, urlHandler)
	go func() {
		logger.Fatal(s.Start())
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	logger.Infof("shutting down by signal")
	if err = stor.Shutdown(); err != nil {
		logger.Error(err)
	}
	os.Exit(0)
}
