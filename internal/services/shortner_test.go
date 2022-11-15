package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/maxsnegir/url-shortener/cmd/config"
	"github.com/maxsnegir/url-shortener/internal/storage"
)

// TestSetURL Проверка того, что данные записываются в хранилище
func TestSetURL(t *testing.T) {
	cfg, _ := config.NewConfig()
	DB := storage.NewURLMapStorage()
	shortener := NewShortener(DB, cfg.Shortener.BaseURL)
	tests := []struct {
		name     string
		url      string
		expected storage.URLData
		userID   string
	}{
		{
			name: "Practikum URL",
			url:  "https://practicum.yandex.ru",
			expected: storage.URLData{
				ShortURL:    "http://localhost:8080/7CwAhsKq/",
				OriginalURL: "https://practicum.yandex.ru",
			},
			userID: "1",
		},
		{
			name: "Stackoverflow URL",
			url:  "https://stackoverflow.com",
			expected: storage.URLData{
				ShortURL:    "http://localhost:8080/UH9SDDvz/",
				OriginalURL: "https://stackoverflow.com",
			},
			userID: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := shortener.SetShortURL(tt.userID, tt.url)
			require.NoError(t, err, "Error while set URL")
			urls, err := DB.GetURLData(tt.userID)
			require.NoError(t, err, "Error while get data from DB")
			require.NotEqual(t, len(urls), 0, "Data in storage is empty")
			assert.Equal(t, tt.expected, urls[len(urls)-1], "unexpected value")
		})
	}
}

// TestGetURLByID Проверка того, что shortener возвращает правильные ссылки по ID
func TestGetURLByID(t *testing.T) {
	DB := storage.NewURLMapStorage()
	shortener := NewShortener(DB, config.BaseURL)
	tests := []struct {
		name     string
		url      string
		expected string
		userID   string
	}{
		{
			name:     "Practikum URL",
			url:      "https://practicum.yandex.ru",
			expected: "https://practicum.yandex.ru",
			userID:   "1",
		},
		{
			name:     "Stackoverflow URL",
			url:      "https://stackoverflow.com",
			expected: "https://stackoverflow.com",
			userID:   "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortURL, err := shortener.SetShortURL(tt.userID, tt.url)
			require.NoError(t, err, "Error while setting URL")
			//urlID := shortener.getURLIdFromShortURL(shortURL)
			originalURL, err := shortener.GetOriginalURLByShort(tt.userID, shortURL)
			require.NoError(t, err, "Error while getting original URL")
			assert.Equal(t, originalURL, tt.expected, "GetOriginalURLByShort return wrong data")
		})
	}
}

// TestParseURL Проверка того, что правильно проверяется валидность URL
func TestParseURL(t *testing.T) {
	DB := storage.NewURLMapStorage()
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
			err := shortener.isURLValid(tt.value)
			if err != nil {
				assert.True(t, tt.wantError, "Unexpected error")
				assert.ErrorIs(t, err, URLIsNotValidError{tt.value}, "Wrong error type")
			}
		})
	}
}
