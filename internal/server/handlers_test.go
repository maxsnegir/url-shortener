package server

import (
	"bytes"
	"encoding/json"
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

func TestSetURLTextHandler(t *testing.T) {
	cfg, _ := config.NewConfig()
	shortURLAddress := cfg.Shortener.BaseURL
	urlDB := storages.NewURLDataBase()
	shortener := services.NewShortener(urlDB, shortURLAddress)
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
				response:    fmt.Sprintf("%s/pkmdI_i-/", shortURLAddress),
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
			router.HandleFunc("/", handler.SetURLTextHandler()).Methods(http.MethodPost)
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

func TestSetURLJSONHandler(t *testing.T) {
	cfg, _ := config.NewConfig()
	shortURLAddress := cfg.Shortener.BaseURL
	urlDB := storages.NewURLDataBase()
	shortener := services.NewShortener(urlDB, shortURLAddress)
	handler := NewURLHandler(shortener, logrus.New())
	type ResponseData struct {
		Result string `json:"result"`
		ErrMsg string `json:"error,omitempty"`
	}

	type want struct {
		code            int
		hasResponseBody bool
		response        ResponseData
		contentType     string
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
				hasResponseBody: true,
				code:            http.StatusCreated,
				response:        ResponseData{Result: fmt.Sprintf("%s/pkmdI_i-/", shortURLAddress)},
				contentType:     "application/json",
			},
			url:    "https://practicum.yandex.ru/",
			method: http.MethodPost,
		},
		{
			name: "Wrong HTTP Method",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
			url:    "https://practicum.yandex.ru/",
			method: http.MethodGet,
		},
		{
			name: "Wrong Body",
			want: want{
				code:            http.StatusBadRequest,
				hasResponseBody: true,
				response: ResponseData{
					ErrMsg: services.URLIsNotValidError{URL: "URL"}.Error(),
				},
				contentType: "application/json",
			},
			url:    "URL",
			method: http.MethodPost,
		},
		{
			name: "Empty Body",
			want: want{
				code:            http.StatusBadRequest,
				hasResponseBody: true,
				response:        ResponseData{ErrMsg: "Wrong URL type or URL in body is missing"},
				contentType:     "application/json",
			},
			url:    "",
			method: http.MethodPost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := httptest.NewRecorder()
			var jsonStr = []byte(fmt.Sprintf(`{"url":"%s"}`, tt.url))
			request := httptest.NewRequest(tt.method, "/api/shorten", bytes.NewBuffer(jsonStr))
			router := mux.NewRouter()
			router.HandleFunc("/api/shorten", handler.SetURLJSONHandler()).Methods(http.MethodPost)
			router.ServeHTTP(writer, request)
			response := writer.Result()
			defer response.Body.Close()

			responseData := &ResponseData{}
			assert.Equal(t, tt.want.code, response.StatusCode, "wrong status code")
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"), "wrong content type")
			if tt.want.hasResponseBody {
				_ = json.NewDecoder(response.Body).Decode(responseData)
				assert.Equal(t, tt.want.response.Result, responseData.Result, "wrong url in response")
				assert.Equal(t, tt.want.response.ErrMsg, responseData.ErrMsg, "wrong error message in response")
			}
		},
		)
	}

}

func TestGetURLByIDHandler(t *testing.T) {
	cfg, _ := config.NewConfig()
	shortURLAddress := cfg.Shortener.BaseURL
	urlDB := storages.NewURLDataBase()
	shortener := services.NewShortener(urlDB, shortURLAddress)
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
