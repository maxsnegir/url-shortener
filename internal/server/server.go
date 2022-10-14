package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/sirupsen/logrus"
	"net/http"
)

type server struct {
	config     config.Config
	logger     *logrus.Logger
	router     *mux.Router
	urlHandler URLHandler
}

func (s *server) Start() error {
	s.beforeStart()
	serverAddress := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	s.logger.Infof("Server is starting on %s", serverAddress)
	return http.ListenAndServe(serverAddress, s.router)
}

func (s *server) beforeStart() {
	s.configureRouter()
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/", s.urlHandler.SetURLHandler()).Methods(http.MethodPost)
	s.router.HandleFunc("/{urlID}/", s.urlHandler.GetURLByIDHandler()).Methods(http.MethodGet)
	s.router.Use(s.LoggingMiddleware)
}

func NewServer(cfg config.Config, logger *logrus.Logger, urlHandler URLHandler) *server {
	return &server{
		router:     mux.NewRouter(),
		config:     cfg,
		logger:     logger,
		urlHandler: urlHandler,
	}
}
