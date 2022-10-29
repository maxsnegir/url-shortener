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
	DB := storage.NewMapURLDataBase()
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
			shortURL, err := shortener.SetURL(tt.value)
			require.NoError(t, err, "Error while set URL")
			urlID := shortener.getURLIdFromShortURL(shortURL)
			value, err := DB.Get(urlID)
			require.NoError(t, err, "Error while get data from DB")
			assert.Equal(t, tt.expected, value, "unexpected value")
		})
	}
}

// TestGetURLByID Проверка того, что shortener возвращает правильные ссылки по ID
func TestGetURLByID(t *testing.T) {
	DB := storage.NewMapURLDataBase()
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
			shortURL, err := shortener.SetURL(tt.value)
			require.NoError(t, err, "Error while setting URL")
			urlID := shortener.getURLIdFromShortURL(shortURL)
			originalURL, err := shortener.GetURLByID(urlID)
			require.NoError(t, err, "Error while getting original URL")
			assert.Equal(t, originalURL, tt.expected, "GetURLByID return wrong data")
		})
	}
}

// TestParseURL Проверка того, что правильно проверяется валидность URL
func TestParseURL(t *testing.T) {
	DB := storage.NewMapURLDataBase()
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
