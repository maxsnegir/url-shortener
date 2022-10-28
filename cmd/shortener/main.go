package main

import (
	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/logging"
	"github.com/maxsnegir/url-shortener/internal/server"
	"github.com/maxsnegir/url-shortener/internal/services"
	"github.com/maxsnegir/url-shortener/internal/storages"
	"log"
	"os"
	"os/signal"
)

func getStorage(cfg config.Config) (storages.Storage, error) {
	switch cfg.Shortener.FileStoragePath {
	case "":
		return storages.NewMapURLDataBase(), nil
	default:
		return storages.NewFileStorage(cfg.Shortener.FileStoragePath)
	}
}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	logger := logging.NewLogger(cfg.Logger.LogLevel)
	storage, err := getStorage(cfg)
	if err != nil {
		switch err.(type) {
		// Если случилась ошибка при загрузке данных из файла - все равно продолжим работу
		case storages.LoadingDumbDataError:
			logger.Error(err)
		default:
			logger.Fatal(err)
		}
	}
	shortener := services.NewShortener(storage, cfg.Shortener.BaseURL)
	urlHandler := server.NewURLHandler(shortener, logger)
	s := server.NewServer(cfg, logger, urlHandler)
	go func() {
		logger.Fatal(s.Start())
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	logger.Infof("shutting down by signal")
	os.Exit(0)
}
