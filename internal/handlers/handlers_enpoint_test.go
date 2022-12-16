package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/auth"
	"github.com/maxsnegir/url-shortener/internal/mocks"
	"github.com/maxsnegir/url-shortener/internal/services"
	"github.com/maxsnegir/url-shortener/internal/storage"
)

func TestSetURLTextHandler(t *testing.T) {
	shortURLAddress := config.BaseURL
	urlDB := storage.NewMemoryURLStorage(storage.NewMapStorage())
	shortener := services.NewShortener(urlDB, shortURLAddress)
	authorization, _ := auth.NewCookieAuthentication("secretKey")
	handler := NewURLHandler(shortener, authorization, logrus.New())
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
			router.Use(handler.CookieAuthenticationMiddleware)
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
	shortURLAddress := config.BaseURL
	urlDB := storage.NewMemoryURLStorage(storage.NewMapStorage())
	shortener := services.NewShortener(urlDB, shortURLAddress)
	authorization, _ := auth.NewCookieAuthentication("secretKey")
	handler := NewURLHandler(shortener, authorization, logrus.New())
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
			router.Use(handler.CookieAuthenticationMiddleware)
			router.ServeHTTP(writer, request)
			response := writer.Result()
			defer response.Body.Close()

			responseData := &ResponseData{}
			assert.Equal(t, tt.want.code, response.StatusCode, "wrong status code")
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"), "wrong content type")
			if tt.want.hasResponseBody {
				require.NoError(t, json.NewDecoder(response.Body).Decode(responseData))
				assert.Equal(t, tt.want.response.Result, responseData.Result, "wrong url in response")
				assert.Equal(t, tt.want.response.ErrMsg, responseData.ErrMsg, "wrong error message in response")
			}
		},
		)
	}

}

func TestGetURLByIDHandler(t *testing.T) {
	shortURLAddress := config.BaseURL
	urlDB := storage.NewMemoryURLStorage(storage.NewMapStorage())
	shortener := services.NewShortener(urlDB, shortURLAddress)
	authorization, _ := auth.NewCookieAuthentication("secretKey")
	handler := NewURLHandler(shortener, authorization, logrus.New())
	type want struct {
		code        int
		response    string
		contentType string
		location    string
	}

	tests := []struct {
		name      string
		want      want
		url       string
		userToken string
		method    string
	}{
		{
			name: "All correct",
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    "",
				contentType: "text/plain; charset=utf-8",
				location:    "https://practicum.yandex.ru/",
			},
			url:       "https://practicum.yandex.ru/",
			userToken: "123",
			method:    http.MethodGet,
		},
		{
			name: "Wrong HTTP Method",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
			url:       "https://github.com/",
			userToken: "123",
			method:    http.MethodPost,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/{urlID}/", handler.GetURLByIDHandler()).Methods(http.MethodGet)
			shortURL, err := shortener.SaveData(context.Background(), tt.userToken, tt.url)
			require.NoError(t, err, "error while saving url")
			request := httptest.NewRequest(tt.method, shortURL, nil)
			router.Use(handler.CookieAuthenticationMiddleware)
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

func TestGetUserURLs(t *testing.T) {
	shortURLAddress := config.BaseURL
	urlDB := storage.NewMemoryURLStorage(storage.NewMapStorage())
	shortener := services.NewShortener(urlDB, shortURLAddress)
	authorization, _ := auth.NewCookieAuthentication("secretKey")
	handler := NewURLHandler(shortener, authorization, logrus.New())
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name   string
		want   want
		url    string
		body   string
		method string
	}{
		{
			name: "All correct",
			body: "https://practicum.yandex.ru/",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "All correct",
			body: "https://practicum.yandex.ru/123",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "All correct",
			body: "https://practicum.yandex.ru/321",
			want: want{
				code: http.StatusOK,
			},
		},
	}
	authToken := ""
	for _, tt := range tests {
		w := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
		if authToken != "" {
			request.AddCookie(&http.Cookie{
				Name:  auth.AuthorizationCookieName,
				Value: authToken,
			})
		}
		router := mux.NewRouter()
		router.HandleFunc("/", handler.SetURLTextHandler()).Methods(http.MethodPost)
		router.Use(handler.CookieAuthenticationMiddleware)
		router.ServeHTTP(w, request)

		if authToken == "" {
			response := w.Result()
			defer response.Body.Close()

			for _, cookie := range response.Cookies() {
				if cookie.Name == auth.AuthorizationCookieName {
					authToken = cookie.Value
				}
			}
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			request := httptest.NewRequest(tt.method, "/api/user/urls", nil)
			request.AddCookie(&http.Cookie{
				Name:  auth.AuthorizationCookieName,
				Value: authToken,
			})
			router := mux.NewRouter()
			router.HandleFunc("/api/user/urls", handler.GetUserURLs()).Methods(http.MethodGet)
			router.Use(handler.CookieAuthenticationMiddleware)
			router.ServeHTTP(w, request)
			response := w.Result()
			defer response.Body.Close()
			var responseData []storage.URLData
			require.NoError(t, json.NewDecoder(response.Body).Decode(&responseData))
			require.Equal(t, len(tests), len(responseData))

		})
	}
}

func TestPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mocks.NewMockShortenerStorage(ctrl)
	shortener := services.NewShortener(s, config.BaseURL)
	authorization, _ := auth.NewCookieAuthentication("secretKey")
	handler := NewURLHandler(shortener, authorization, logrus.New())

	tests := []struct {
		name       string
		statusCode int
		err        error
	}{
		{
			name:       "All ok",
			statusCode: http.StatusOK,
			err:        nil,
		},
		{
			name:       "Want error",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.EXPECT().Ping(gomock.Any()).Return(tt.err)

			request := httptest.NewRequest(http.MethodGet, "/ping", nil)
			router := mux.NewRouter()
			router.HandleFunc("/ping", handler.Ping()).Methods(http.MethodGet)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			response := w.Result()
			defer response.Body.Close()
			require.Equal(t, response.StatusCode, tt.statusCode, "wrong status code")
		})
	}
}

func TestDuplicateError(t *testing.T) {
	tests := []struct {
		name         string
		originalURL  string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "First save no err",
			originalURL:  "https://github.com",
			expectedCode: http.StatusCreated,
			expectedBody: "http://localhost:8080/hLfkSqVN/",
		},
		{
			name:         "Second save conflict err",
			originalURL:  "https://github.com",
			expectedCode: http.StatusConflict,
			expectedBody: "http://localhost:8080/hLfkSqVN/",
		},
	}
	shortener := services.NewShortener(storage.NewMemoryURLStorage(storage.NewMapStorage()), config.BaseURL)
	authorization, _ := auth.NewCookieAuthentication("secretKey")
	handler := NewURLHandler(shortener, authorization, logrus.New())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.originalURL))
			router := mux.NewRouter()
			router.HandleFunc("/", handler.SetURLTextHandler()).Methods(http.MethodPost)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			response := w.Result()
			defer response.Body.Close()
			require.Equal(t, tt.expectedCode, response.StatusCode, "wrong status code")
			resBody, err := io.ReadAll(response.Body)
			require.NoError(t, err, "error while reading response body")
			require.Equal(t, tt.expectedBody, string(resBody), "wrong response body")
		})
	}
}

func TestDeleteURLs(t *testing.T) {
	shortener := services.NewShortener(storage.NewMemoryURLStorage(storage.NewMapStorage()), config.BaseURL)
	authorization, _ := auth.NewCookieAuthentication("secretKey")
	handler := NewURLHandler(shortener, authorization, logrus.New())

	testURLs := map[string]string{
		"http://github.com/":            "",
		"http://gitlab.com":             "",
		"https://bitbucket.org":         "",
		"https://www.mercurial-scm.org": "",
	}

	body := make([]string, 0, len(testURLs))
	t.Run("create short urls", func(t *testing.T) {
		for fullURL := range testURLs {
			shortURL, err := shortener.SaveData(context.Background(), "someToken", fullURL)
			require.NoError(t, err, "error while create short url")
			testURLs[fullURL] = shortURL
			body = append(body, shortener.GetURLIdFromShortURL(shortURL))
		}
	})

	t.Run("delete urls", func(t *testing.T) {
		b, err := json.Marshal(body)
		require.NoError(t, err, "error while encoding request body")
		request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBuffer(b))
		router := mux.NewRouter()
		router.HandleFunc("/api/user/urls", handler.DeleteURLS()).Methods(http.MethodDelete)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		response := w.Result()
		defer response.Body.Close()
		require.Equal(t, http.StatusAccepted, response.StatusCode, "wrong status code")
	})

	t.Run("check url is deleted", func(t *testing.T) {
		checkURLIsDeleted := func(shortURL string) bool {
			w := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/{urlID}/", handler.GetURLByIDHandler()).Methods(http.MethodGet)
			request := httptest.NewRequest(http.MethodGet, shortURL, nil)
			router.Use(handler.CookieAuthenticationMiddleware)
			router.ServeHTTP(w, request)
			response := w.Result()
			defer response.Body.Close()
			return response.StatusCode == http.StatusGone
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				t.Error("not all urls are deleted")
				return
			default:
				deletedURLsCount := 0
				for _, shortURL := range testURLs {
					if checkURLIsDeleted(shortURL) {
						deletedURLsCount++
					}
				}
				if deletedURLsCount == len(testURLs) {
					return
				}
			}
		}
	})
}
