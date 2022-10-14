package server

import (
	"github.com/gorilla/mux"
	"github.com/maxsnegir/url-shortener/internal/services"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type BaseHandler struct {
	logger *logrus.Logger
}

func (h *BaseHandler) TextResponse(w http.ResponseWriter, code int, data string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)

	if data != "" {
		if _, err := w.Write([]byte(data)); err != nil {
			h.logger.Error(err)
		}
	}
}

type URLHandler struct {
	BaseHandler
	shortener services.URLService
}

func (h *URLHandler) SetURLHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		url, err := io.ReadAll(r.Body)
		if len(url) == 0 || err != nil {
			h.TextResponse(w, http.StatusUnprocessableEntity, "URL in request body is missing")
			return
		}
		shortURL, err := h.shortener.SetURL(string(url))
		if err != nil {
			switch err.(type) {
			case services.URLIsNotValidError:
				h.TextResponse(w, http.StatusBadRequest, err.Error())
			default:
				h.logger.Error(err)
				h.TextResponse(w, http.StatusInternalServerError, InternalServerError.Error())
			}
			return
		}
		h.TextResponse(w, http.StatusCreated, shortURL)
	}
}

func (h *URLHandler) GetURLByIDHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		urlID := vars["urlID"]
		originalURL, err := h.shortener.GetURLByID(urlID)
		if err != nil {
			switch err.(type) {
			case services.OriginalURLNotFound:
				h.TextResponse(w, http.StatusNotFound, err.Error())
			default:
				h.logger.Error(err)
				h.TextResponse(w, http.StatusInternalServerError, InternalServerError.Error())
			}
			return
		}
		w.Header().Add("Location", originalURL)
		h.TextResponse(w, http.StatusTemporaryRedirect, "")
	}
}

func NewURLHandler(shortener services.URLService, logger *logrus.Logger) URLHandler {
	return URLHandler{
		BaseHandler: BaseHandler{logger: logger},
		shortener:   shortener,
	}
}
