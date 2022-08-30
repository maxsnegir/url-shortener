package server

import (
	"encoding/json"
	"fmt"
	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/databases"
	"github.com/maxsnegir/url-shortener/internal/services"
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

func (s *server) Response(w http.ResponseWriter, r *http.Request, code int, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err != nil {
		data = map[string]string{"error": err.Error()}
	}
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			code = http.StatusInternalServerError
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

func NewServer(cfg config.Config, logger *logrus.Logger, db databases.KeyValueDB) *server {
	shortener := services.NewShortener(db)
	return &server{
		Config:        cfg,
		Logger:        logger,
		Shortener:     shortener,
		DefaultRouter: NotFoundHandler(),
	}
}
