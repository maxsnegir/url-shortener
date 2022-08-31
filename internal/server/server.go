package server

import (
	"fmt"
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
	Routes        []Router
	DefaultRouter http.HandlerFunc
	// Сервисы приложения
	Shortener *services.Shortener
}

func (s *server) configureRouter() {
	s.Handler(`^/$`, s.SetURLHandler())
	s.Handler(`^/.*?/$`, s.GetURLByIDHandler())
}

func (s *server) TextResponse(w http.ResponseWriter, r *http.Request, code int, data string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)

	if data != "" {
		if _, err := w.Write([]byte(data)); err != nil {
			s.Logger.Error(err)
		}
	}
	s.Logger.Infof("%s::%s::%v", r.RequestURI, r.Method, code)
}

func (s *server) Start() error {
	s.configureRouter()
	serverAddress := fmt.Sprintf("%s:%s", s.Config.Server.Host, s.Config.Server.Port)
	s.Logger.Infof("Server is starting on %s", serverAddress)
	err := http.ListenAndServe(serverAddress, s)
	if err != nil {
		return err
	}
	return nil
}

func NewServer(cfg config.Config, logger *logrus.Logger, db storages.KeyValueStorage) *server {
	shortener := services.NewShortener(db)
	return &server{
		Config:        cfg,
		Logger:        logger,
		Shortener:     shortener,
		DefaultRouter: NotFoundHandler(),
	}
}
