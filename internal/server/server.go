package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/services"
	"github.com/maxsnegir/url-shortener/internal/storages"
	"github.com/sirupsen/logrus"
	"net/http"
)

type server struct {
	// Хелперы
	Config config.Config
	Logger *logrus.Logger
	// Роутинг
	Router *mux.Router
	// Сервисы приложения
	Shortener *services.Shortener
}

func (s *server) configureRouter() {
	s.Router.HandleFunc("/", s.SetURLHandler())
	s.Router.HandleFunc("/{urlID}/", s.GetURLByIDHandler())
	s.Router.Use(s.LoggingMiddleware)
}

func (s *server) Start() error {
	s.configureRouter()
	serverAddress := fmt.Sprintf("%s:%s", s.Config.Server.Host, s.Config.Server.Port)
	s.Logger.Infof("Server is starting on %s", serverAddress)
	err := http.ListenAndServe(serverAddress, s.Router)
	if err != nil {
		return err
	}
	return nil
}

func NewServer(cfg config.Config, logger *logrus.Logger, db storages.KeyValueStorage) *server {
	shortener := services.NewShortener(db)
	return &server{
		Router:    mux.NewRouter(),
		Config:    cfg,
		Logger:    logger,
		Shortener: shortener,
	}
}
