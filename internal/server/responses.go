package server

import "net/http"

func (s *server) TextResponse(w http.ResponseWriter, code int, data string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)

	if data != "" {
		if _, err := w.Write([]byte(data)); err != nil {
			s.Logger.Error(err)
		}
	}
}
