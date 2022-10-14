package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/services"
	"github.com/maxsnegir/url-shortener/internal/storages"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSetURLHandler(t *testing.T) {
	cfg := config.NewConfig("../../configs/config.yaml")
	serverAddress := cfg.Server.FullAddress
	urlDB := storages.NewURLDataBase()
	shortener := services.NewShortener(urlDB, serverAddress)
	handler := NewURLHandler(shortener, logrus.New())
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
				code: http.StatusMethodNotAllowed,
			},
			body:   "https://practicum.yandex.ru/",
			method: http.MethodGet,
		},
		{
			name: "Wrong Body",
			want: want{
				code:        http.StatusBadRequest,
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
			w := httptest.NewRecorder()

			request := httptest.NewRequest(tt.method, "/", body)
			router := mux.NewRouter()
			router.HandleFunc("/", handler.SetURLHandler()).Methods(http.MethodPost)
			router.ServeHTTP(w, request)

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
	serverAddress := cfg.Server.FullAddress
	urlDB := storages.NewURLDataBase()
	shortener := services.NewShortener(urlDB, serverAddress)
	handler := NewURLHandler(shortener, logrus.New())
	type want struct {
		code        int
		response    string
		contentType string
		location    string
	}

	tests := []struct {
		name   string
		want   want
		url    string
		method string
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
				code: http.StatusMethodNotAllowed,
			},
			url:    "https://practicum.yandex.ru/",
			method: http.MethodPost,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/{urlID}/", handler.GetURLByIDHandler()).Methods(http.MethodGet)
			shortURL, _ := shortener.SetURL(tt.url)
			request := httptest.NewRequest(tt.method, shortURL, nil)
			router.ServeHTTP(w, request)
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
