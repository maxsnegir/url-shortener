package storage

import (
	"context"

	"github.com/maxsnegir/url-shortener/cmd/config"
)

type URLData struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Storage interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Shutdown(ctx context.Context) error
}

type ShortenerStorage interface {
	GetOriginalURL(shortURL string) (string, error)
	SetShortURL(urlData URLData) error
	GetUserURLs(userID string) ([]string, error)
	SetUserURL(userID string, shortURL string) error
	Shutdown(ctx context.Context) error
	Ping(ctx context.Context) error
}

func GetURLStorage(cfg config.Config) (ShortenerStorage, error) {
	//PostgresStorage
	if cfg.Storage.DatabaseDSN != "" {
		return NewPostgresStorage(context.Background(), cfg.Storage.DatabaseDSN)
	}
	//MapStorage
	if cfg.Storage.FileStoragePath == "" {
		return NewURLStorage(NewMapStorage()), nil
	}
	//FileStorage
	fileStorage, err := NewURLFileStorage(cfg.Storage.FileStoragePath)
	if err != nil {
		return nil, err
	}
	return NewURLStorage(fileStorage), nil

}
