package server

import (
	"net/http"
)

func (s *server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		s.logger.Infof("%s :: %s", r.RequestURI, r.Method)
	})
}
