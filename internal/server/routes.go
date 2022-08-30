package server

import (
	"net/http"
	"regexp"
)

type Router struct {
	Pattern     *regexp.Regexp
	Handler     http.HandlerFunc
	Middlewares []Middleware
}

func (s *server) Handler(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	re := regexp.MustCompile(pattern)
	route := Router{Pattern: re, Handler: handler, Middlewares: middlewares}
	s.Routes = append(s.Routes, route)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	for _, rt := range s.Routes {
		if matches := rt.Pattern.FindStringSubmatch(r.URL.Path); len(matches) > 0 {
			handler := MiddlewareConveyor(rt.Handler, rt.Middlewares...)
			handler(w, r)
			return
		}
	}
	s.DefaultRouter(w, r)
}
