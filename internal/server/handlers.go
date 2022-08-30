package server

import (
	"github.com/maxsnegir/url-shortener/internal/services"
	"io"
	"net/http"
	"strings"
)

func (s *server) SetURLHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s.Response(w, r, http.StatusMethodNotAllowed, nil, MethodNotAllowedError{r.Method})
			return
		}

		url, err := io.ReadAll(r.Body)
		if len(url) == 0 || err != nil {
			s.Response(w, r, http.StatusUnprocessableEntity, nil, RequestParamsError{})
			return
		}
		stringURL := string(url)
		if !s.Shortener.URLIsValid(stringURL) {
			s.Response(w, r, http.StatusUnprocessableEntity, nil, services.URLIsNotValidError{Url: stringURL})
			return
		}
		// Пока нет тз, пусть ссылка хранится вечно(подразумевалась для редиса)
		shortURLID, err := s.Shortener.SetURL(stringURL, 0)
		if err != nil {
			s.Logger.Error(err)
			s.Response(w, r, http.StatusInternalServerError, nil, InternalServerError{})
			return
		}
		s.Response(w, r, http.StatusOK, shortURLID, nil)
	}
}

func (s *server) GetURLByIDHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			s.Response(w, r, http.StatusMethodNotAllowed, nil, MethodNotAllowedError{r.Method})
			return
		}
		urlId := strings.Split(r.URL.Path, "/")[1]
		originalUrl, err := s.Shortener.GetURLByID(urlId)
		if err != nil {
			switch err.(type) {
			case services.OriginalURLNotFound:
				s.Response(w, r, http.StatusNotFound, nil, err)
			default:
				s.Logger.Error(err)
				s.Response(w, r, http.StatusInternalServerError, nil, InternalServerError{})
			}
			return
		}
		w.Header().Add("Location", originalUrl)
		s.Response(w, r, http.StatusTemporaryRedirect, nil, nil)
	}
}

func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, NotFoundError{r.URL.String()}.Error(), http.StatusNotFound)
	}
}
