package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/auth"
	"github.com/maxsnegir/url-shortener/internal/services"
	"github.com/maxsnegir/url-shortener/internal/storage"
)

// TestCookieAuthMiddleware проверяем наличие возвращаемой куки авторизации
func TestCookieAuthMiddleware(t *testing.T) {
	urlStorage := storage.NewURLStorage(storage.NewMapStorage())
	shortener := services.NewShortener(urlStorage, config.BaseURL)
	authorization, _ := auth.NewCookieAuthentication("secret")
	handler := NewURLHandler(shortener, authorization, logrus.New())
	tests := []struct {
		name       string
		path       string
		method     string
		handleFunc func() http.HandlerFunc
	}{
		{
			name:       "Set text short url",
			path:       "/",
			method:     http.MethodPost,
			handleFunc: handler.SetURLTextHandler,
		},
		{
			name:       "Set json short url",
			path:       "/api/shorten",
			method:     http.MethodPost,
			handleFunc: handler.SetURLJSONHandler,
		},
		{
			name:       "Get url by short url",
			path:       "/something",
			method:     http.MethodGet,
			handleFunc: handler.GetURLByIDHandler,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			request := httptest.NewRequest(tt.method, "/", nil)
			router := mux.NewRouter()
			router.HandleFunc("/", tt.handleFunc()).Methods(tt.method)
			router.Use(handler.CookieAuthenticationMiddleware)
			router.ServeHTTP(w, request)
			response := w.Result()
			defer response.Body.Close()
			var authCookie string
			for _, cookie := range response.Cookies() {
				if cookie.Name == auth.AuthorizationCookieName {
					authCookie = cookie.Value
				}
			}
			require.NotEmpty(t, authCookie, "Auth Cookie: is empty")
			userToken, err := authorization.ParseToken(authCookie)
			require.NoError(t, err, "Error while parsing auth cookie")
			require.NotEmpty(t, userToken, "User Token is empty")
		})
	}
}
