package server

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func MiddlewareConveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}
