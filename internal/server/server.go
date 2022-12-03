package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/handlers"
)

type server struct {
	config     config.Config
	logger     *logrus.Logger
	router     *mux.Router
	urlHandler handlers.URLHandler
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
	s.router.HandleFunc("/api/user/urls", s.urlHandler.GetUserURLs()).Methods(http.MethodGet)
	s.router.HandleFunc("/ping", s.urlHandler.Ping()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/shorten/batch", s.urlHandler.SaveDataBatch()).Methods(http.MethodPost)
	s.router.Use(s.urlHandler.CookieAuthenticationMiddleware)
	s.router.Use(s.urlHandler.LoggingMiddleware)
	s.router.Use(s.urlHandler.GzipMiddleware)
	s.router.Use(s.urlHandler.UnzipMiddleware)
}

func NewServer(cfg config.Config, logger *logrus.Logger, urlHandler handlers.URLHandler) *server {
	return &server{
		router:     mux.NewRouter(),
		config:     cfg,
		logger:     logger,
		urlHandler: urlHandler,
	}
}
