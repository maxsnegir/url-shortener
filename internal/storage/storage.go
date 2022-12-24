package storage

import (
	"context"

	"github.com/maxsnegir/url-shortener/cmd/config"
)

type URLData struct {
	URLDataID   int    `json:"-" db:"url_data_id"`
	ShortURL    string `json:"short_url" db:"short_url"`
	OriginalURL string `json:"original_url" db:"original_url"`
	Deleted     bool   `json:"deleted" db:"deleted"`
}

type Storage interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Shutdown() error
}

type ShortenerStorage interface {
	SaveData(ctx context.Context, userToken string, urlData URLData) error
	SaveDataBatch(ctx context.Context, userToken string, urlData []URLData) error
	GetOriginalURL(ctx context.Context, shortURL string) (URLData, error)
	GetUserURLs(ctx context.Context, userToken string) ([]URLData, error)
	Shutdown() error
	Ping(ctx context.Context) error
	DeleteURLs(ctx context.Context, urlsToDelete []string) error
}

func GetURLStorage(cfg config.Config) (ShortenerStorage, error) {
	//PostgresStorage
	if cfg.Storage.DatabaseDSN != "" {
		return NewPostgresStorage(context.Background(), cfg.Storage.DatabaseDSN)
	}
	//MapStorage
	if cfg.Storage.FileStoragePath == "" {
		return NewMemoryURLStorage(NewMapStorage()), nil
	}
	//FileStorage
	fileStorage, err := NewURLFileStorage(cfg.Storage.FileStoragePath)
	if err != nil {
		return nil, err
	}
	return NewMemoryURLStorage(fileStorage), nil
}
