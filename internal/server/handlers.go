package server

import (
	"fmt"
	"github.com/maxsnegir/url-shortener/internal/services"
	"io"
	"net/http"
	"strings"
)

func (s *server) SetURLHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s.TextResponse(w, r, http.StatusMethodNotAllowed, MethodNotAllowedError{r.Method}.Error())
			return
		}

		url, err := io.ReadAll(r.Body)
		if len(url) == 0 || err != nil {
			s.TextResponse(w, r, http.StatusUnprocessableEntity, "URL in request body is missing")
			return
		}
		originalURL, err := s.Shortener.ParseURL(string(url))
		if err != nil {
			s.TextResponse(w, r, http.StatusUnprocessableEntity, err.Error())
			return
		}
		// Пока нет тз, пусть ссылка хранится вечно(подразумевалась для редиса)
		urlID, err := s.Shortener.SetURL(originalURL, 0)
		if err != nil {
			s.Logger.Error(err)
			s.TextResponse(w, r, http.StatusInternalServerError, InternalServerError.Error())
			return
		}
		shortURL := fmt.Sprintf("%s/%s/", s.Config.Server.FullAddress, urlID)
		s.TextResponse(w, r, http.StatusCreated, shortURL)
	}
}

func (s *server) GetURLByIDHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			s.TextResponse(w, r, http.StatusMethodNotAllowed, MethodNotAllowedError{r.Method}.Error())
			return
		}
		urlID := strings.Split(r.URL.Path, "/")[1]
		originalURL, err := s.Shortener.GetURLByID(urlID)
		if err != nil {
			switch err.(type) {
			case services.OriginalURLNotFound:
				s.TextResponse(w, r, http.StatusNotFound, err.Error())
			default:
				s.Logger.Error(err)
				s.TextResponse(w, r, http.StatusInternalServerError, InternalServerError.Error())
			}
			return
		}
		w.Header().Add("Location", originalURL)
		s.TextResponse(w, r, http.StatusTemporaryRedirect, "")
	}
}

func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, NotFoundError{r.URL.String()}.Error(), http.StatusNotFound)
	}
}
