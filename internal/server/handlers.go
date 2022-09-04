package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/maxsnegir/url-shortener/internal/services"
	"io"
	"net/http"
)

func (s *server) SetURLHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s.TextResponse(w, http.StatusMethodNotAllowed, MethodNotAllowedError{r.Method}.Error())
			return
		}

		url, err := io.ReadAll(r.Body)
		if len(url) == 0 || err != nil {
			s.TextResponse(w, http.StatusUnprocessableEntity, "URL in request body is missing")
			return
		}
		originalURL, err := s.Shortener.ParseURL(string(url))
		if err != nil {
			s.TextResponse(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		// Пока нет тз, пусть ссылка хранится вечно(логика подразумевалась для хранения ссылки в редисе)
		urlID, err := s.Shortener.SetURL(originalURL, 0)
		if err != nil {
			s.Logger.Error(err)
			s.TextResponse(w, http.StatusInternalServerError, InternalServerError.Error())
			return
		}
		shortURL := fmt.Sprintf("%s/%s/", s.Config.Server.FullAddress, urlID)
		s.TextResponse(w, http.StatusCreated, shortURL)
	}
}

func (s *server) GetURLByIDHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			s.TextResponse(w, http.StatusMethodNotAllowed, MethodNotAllowedError{r.Method}.Error())
			return
		}
		vars := mux.Vars(r)
		urlID := vars["urlID"]
		originalURL, err := s.Shortener.GetURLByID(urlID)
		if err != nil {
			switch err.(type) {
			case services.OriginalURLNotFound:
				s.TextResponse(w, http.StatusNotFound, err.Error())
			default:
				s.Logger.Error(err)
				s.TextResponse(w, http.StatusInternalServerError, InternalServerError.Error())
			}
			return
		}
		w.Header().Add("Location", originalURL)
		s.TextResponse(w, http.StatusTemporaryRedirect, "")
	}
}
