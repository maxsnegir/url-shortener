package server

import (
	"fmt"
	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/sirupsen/logrus"
	"net/http"
)

type server struct {
	config config.Config
	Logger *logrus.Logger
	router *http.ServeMux
}

func (s *server) Start() error {
	serverAdd := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	s.Logger.Infof("Server is starting on %s", serverAdd)
	err := http.ListenAndServe(serverAdd, s.router)
	if err != nil {
		return err
	}
	return nil
}

func NewServer(cfg config.Config, logger *logrus.Logger) *server {

	return &server{
		config: cfg,
		Logger: logger,
		router: http.NewServeMux(),
	}
}
