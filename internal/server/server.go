package server

import (
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
	s.logger.Infof("Server is starting on %s", s.config.Server.ServerAddress)
	return http.ListenAndServe(s.config.Server.ServerAddress, s.router)
}

func (s *server) beforeStart() {
	s.configureRouter()
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/", s.urlHandler.SetURLTextHandler()).Methods(http.MethodPost)
	s.router.HandleFunc("/api/shorten", s.urlHandler.SetURLJSONHandler()).Methods(http.MethodPost)
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
