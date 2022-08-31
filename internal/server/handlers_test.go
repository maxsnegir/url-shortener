package server

import (
	"fmt"
	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/services"
	"github.com/maxsnegir/url-shortener/internal/storages"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSetURLHandler(t *testing.T) {
	cfg := config.NewConfig("../../configs/config.yaml")
	s := NewServer(cfg, logrus.New(), storages.NewURLDateBase())
	serverAddress := cfg.Server.FullAddress
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name   string
		want   want
		body   string
		method string
	}{
		{
			name: "All correct",
			want: want{
				code:        http.StatusCreated,
				response:    fmt.Sprintf("%s/pkmdI_i-/", serverAddress),
				contentType: "text/plain; charset=utf-8",
			},
			body:   "https://practicum.yandex.ru/",
			method: http.MethodPost,
		},
		{
			name: "Wrong HTTP Method",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    MethodNotAllowedError{http.MethodGet}.Error(),
				contentType: "text/plain; charset=utf-8",
			},
			body:   "https://practicum.yandex.ru/",
			method: http.MethodGet,
		},
		{
			name: "Wrong Body",
			want: want{
				code:        http.StatusUnprocessableEntity,
				response:    services.URLIsNotValidError{URL: "URL"}.Error(),
				contentType: "text/plain; charset=utf-8",
			},
			body:   "URL",
			method: http.MethodPost,
		},
		{
			name: "Empty Body",
			want: want{
				code:        http.StatusUnprocessableEntity,
				response:    "URL in request body is missing",
				contentType: "text/plain; charset=utf-8",
			},
			body:   "",
			method: http.MethodPost,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)
			request := httptest.NewRequest(tt.method, "/", body)
			h := s.SetURLHandler()
			w := httptest.NewRecorder()
			h.ServeHTTP(w, request)
			response := w.Result()
			defer response.Body.Close()
			resBody, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want.code, response.StatusCode, "wrong status code")
			assert.Equal(t, tt.want.response, string(resBody), "wrong response body")
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"), "wrong content type")
		})
	}
}

func TestGetURLByIDHandler(t *testing.T) {
	cfg := config.NewConfig("../../configs/config.yaml")
	s := NewServer(cfg, logrus.New(), storages.NewURLDateBase())
	type want struct {
		code        int
		response    string
		contentType string
		location    string
	}

	tests := []struct {
		name     string
		want     want
		url      string
		method   string
		unsetURL bool
	}{
		{
			name: "All correct",
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    "",
				contentType: "text/plain; charset=utf-8",
				location:    "https://practicum.yandex.ru/",
			},
			url:    "https://practicum.yandex.ru/",
			method: http.MethodGet,
		},
		{
			name: "Wrong HTTP Method",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    MethodNotAllowedError{http.MethodPost}.Error(),
				contentType: "text/plain; charset=utf-8",
				location:    "",
			},
			url:    "https://practicum.yandex.ru/",
			method: http.MethodPost,
		},
		{
			name: "Un existing URL ID",
			want: want{
				code:        http.StatusNotFound,
				response:    services.OriginalURLNotFound{URLID: "blabla"}.Error(),
				contentType: "text/plain; charset=utf-8",
				location:    "",
			},
			method:   http.MethodGet,
			unsetURL: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			URL, _ := url.Parse(tt.url)
			var urlID string
			switch tt.unsetURL {
			case true:
				urlID = "blabla" // Несуществующий urlID
			default:
				urlID, _ = s.Shortener.SetURL(URL, 0)
			}

			request := httptest.NewRequest(tt.method, fmt.Sprintf("/%s/", urlID), nil)
			h := s.GetURLByIDHandler()
			w := httptest.NewRecorder()
			h.ServeHTTP(w, request)
			response := w.Result()
			defer response.Body.Close()
			resBody, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want.code, response.StatusCode)
			assert.Equal(t, tt.want.response, string(resBody), "wrong response body")
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"), "wrong contentType")
			assert.Equal(t, tt.want.location, response.Header.Get("Location"), "Wrong Location in header")
		})
	}
}
