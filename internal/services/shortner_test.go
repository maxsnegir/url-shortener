package services

import (
	"github.com/maxsnegir/url-shortener/internal/storages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

// TestSetURL Проверка того, что данные записываются в хранилище
func TestSetURL(t *testing.T) {
	DB := storages.NewURLDateBase()
	shortener := NewShortener(DB)
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
			URL, _ := url.Parse(tt.value)
			urlID, err := shortener.SetURL(URL, 0)
			require.NoError(t, err, "Error while set URL")
			value, err := DB.Get(urlID)
			require.NoError(t, err, "Error while get data from DB")
			stringValue, ok := value.(string)
			require.True(t, ok, "Wrong type, expected=string")
			assert.Equal(t, tt.expected, stringValue, "unexpected value")
		})
	}
}

// TestGetURLByID Проверка того, что Shortener возвращает правильные ссылки по ID
func TestGetURLByID(t *testing.T) {
	DB := storages.NewURLDateBase()
	shortener := NewShortener(DB)
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
			URL, _ := url.Parse(tt.value)
			urlID, err := shortener.SetURL(URL, 0)
			require.NoError(t, err, "Error while setting URL")
			originalURL, err := shortener.GetURLByID(urlID)
			require.NoError(t, err, "Error while getting original URL")
			assert.Equal(t, originalURL, tt.expected, "GetURLByID return wrong data")
		})
	}
}

// TestParseURL Проверка того, что правильно проверяется валидность URL
func TestParseURL(t *testing.T) {
	DB := storages.NewURLDateBase()
	shortener := NewShortener(DB)
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
			_, err := shortener.ParseURL(tt.value)
			if err != nil {
				assert.True(t, tt.wantError, "Unexpected error")
				assert.ErrorIs(t, err, URLIsNotValidError{tt.value}, "Wrong error type")
			}
		})
	}
}
