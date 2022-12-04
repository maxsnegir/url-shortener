package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/storage"
)

// TestSetURL Проверка того, что данные записываются в хранилище
func TestSetURL(t *testing.T) {
	cfg, _ := config.NewConfig()
	DB := storage.NewURLStorage(storage.NewMapStorage())
	shortener := NewShortener(DB, cfg.Shortener.BaseURL)
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "Practikum URL",
			value:    "https://practicum.yandex.ru",
			expected: "https://practicum.yandex.ru",
		},
		{
			name:     "Stackoverflow URL",
			value:    "https://stackoverflow.com",
			expected: "https://stackoverflow.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortURL, err := shortener.SaveData(context.Background(), "", tt.value)
			require.NoError(t, err, "Error while set URL")
			value, err := DB.GetOriginalURL(context.Background(), shortURL)
			require.NoError(t, err, "Error while get data from DB")
			assert.Equal(t, tt.expected, value, "unexpected value")
		})
	}
}

// TestGetURLByID Проверка того, что shortener возвращает правильные ссылки по ID
func TestGetURLByID(t *testing.T) {
	DB := storage.NewURLStorage(storage.NewMapStorage())
	shortener := NewShortener(DB, config.BaseURL)
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "Practikum URL",
			value:    "https://practicum.yandex.ru",
			expected: "https://practicum.yandex.ru",
		},
		{
			name:     "Stackoverflow URL",
			value:    "https://stackoverflow.com",
			expected: "https://stackoverflow.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortURL, err := shortener.SaveData(context.Background(), "", tt.value)
			require.NoError(t, err, "Error while setting URL")
			originalURL, err := DB.GetOriginalURL(context.Background(), shortURL)
			require.NoError(t, err, "Error while getting original URL")
			assert.Equal(t, originalURL, tt.expected, "GetURLByID return wrong data")
		})
	}
}

// TestParseURL Проверка того, что правильно проверяется валидность URL
func TestParseURL(t *testing.T) {
	DB := storage.NewURLStorage(storage.NewMapStorage())
	shortener := NewShortener(DB, config.BaseURL)
	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{
			name:      "Correct URL",
			value:     "https://music.yandex.ru/",
			wantError: false,
		},
		{
			name:      "Wrong URL",
			value:     "music.yandex",
			wantError: true,
		},
		{
			name:      "Empty string",
			value:     "",
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := shortener.IsURLValid(tt.value)
			if err != nil {
				assert.True(t, tt.wantError, "Unexpected error")
				assert.ErrorIs(t, err, URLIsNotValidError{tt.value}, "Wrong error type")
			}
		})
	}
}

func TestGetAllUserURLs(t *testing.T) {
	userToken := "someToken"
	testURLs := map[string]string{
		"http://github.com/":            "http://localhost:8080/hnTOWmuz/",
		"http://gitlab.com":             "http://localhost:8080/PSynyReI/",
		"https://bitbucket.org":         "http://localhost:8080/ji7Semk-/",
		"https://www.mercurial-scm.org": "http://localhost:8080/jdR6WcSi/",
	}

	db := storage.NewURLStorage(storage.NewMapStorage())
	shortener := NewShortener(db, config.BaseURL)
	for shortURL := range testURLs {
		_, err := shortener.SaveData(context.Background(), userToken, shortURL)
		assert.NoError(t, err)
	}

	tests := []struct {
		name         string
		userToken    string
		userURLCount int
	}{
		{
			name:         "All ok",
			userToken:    userToken,
			userURLCount: len(testURLs),
		},
		{
			name:         "Not existing user token",
			userToken:    "There is not user with that auth token",
			userURLCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlData, err := shortener.GetUserURLs(context.Background(), tt.userToken)
			require.NoError(t, err, "error while get all user urls")
			require.Equal(t, len(urlData), tt.userURLCount, "wrong url count")

			for _, url := range urlData {
				shortURL, ok := testURLs[url.OriginalURL]
				require.True(t, ok, "not expected url")
				require.Equal(t, shortURL, url.ShortURL)
			}
		})
	}
}

func TestSaveDataBatch(t *testing.T) {
	DB := storage.NewURLStorage(storage.NewMapStorage())
	shortener := NewShortener(DB, config.BaseURL)
	batchRequest := []URLDataBatchRequest{
		{CorrelationID: "1", OriginalURL: "http://github.com/"},
		{CorrelationID: "2", OriginalURL: "http://gitlab.com"},
		{CorrelationID: "3", OriginalURL: "https://bitbucket.org"},
		{CorrelationID: "4", OriginalURL: "https://www.mercurial-scm.org"},
	}

	t.Run("Run SaveDataBatch", func(t *testing.T) {
		batchResponse, err := shortener.SaveDataBatch(context.Background(), "userToken", batchRequest)
		require.NoError(t, err, "error while make batch save")
		require.Equal(t, len(batchRequest), len(batchResponse))

		for _, urlData := range batchResponse {
			originalURL, err := shortener.GetOriginalURL(context.Background(), urlData.ShortURL)
			require.NoError(t, err, "error while getting original url")

			for _, batchReq := range batchRequest {
				if batchReq.OriginalURL == originalURL {
					require.Equal(t, batchReq.CorrelationID, urlData.CorrelationID, "CorrelationID in request and response didn't match")
					break
				}
			}
		}
	})
}
