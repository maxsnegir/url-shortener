package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type BaseHandler struct {
	logger *logrus.Logger
}

func (h *BaseHandler) TextResponse(w http.ResponseWriter, code int, data string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)

	if _, err := w.Write([]byte(data)); err != nil {
		h.logger.Error(err)
	}
}

func (h *BaseHandler) JSONResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error(err)
	}
}
