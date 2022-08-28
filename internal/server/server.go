package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/services"
	"github.com/sirupsen/logrus"
	"net/http"
)

type server struct {
	// Хелперы
	config config.Config
	Logger *logrus.Logger
	router *http.ServeMux
	// Сервисы приложения
	Shortener *services.Shortener
}

func (s *server) Start() error {

	s.configureRouter()
	serverAdd := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	s.Logger.Infof("Server is starting on %s", serverAdd)
	err := http.ListenAndServe(serverAdd, s.router)
	if err != nil {
		return err
	}
	return nil
}

func (s *server) configureRouter() {
	s.router.Handle("/{:id}", MiddlewareConveyor(s.GetUrlByIdHandler()))
	s.router.Handle("/", MiddlewareConveyor(s.ShortUrlHandler()))
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

func NewServer(cfg config.Config, logger *logrus.Logger, redisClient *redis.Client) *server {
	shortener := services.NewShortener(redisClient)
	return &server{
		config:    cfg,
		Logger:    logger,
		Shortener: shortener,
		router:    http.NewServeMux(),
	}
}
