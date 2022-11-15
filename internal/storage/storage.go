package storage

import "github.com/maxsnegir/url-shortener/cmd/config"

type URLData struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Storage interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Shutdown() error
}

type ShortenerStorage interface {
	GetOriginalURL(shortURL string) (string, error)
	SetShortURL(urlData URLData) error
	GetUserURLs(userID string) ([]string, error)
	SetUserURL(userID string, shortURL string) error
	Shutdown() error
}

func GetURLStorage(cfg config.Config) (ShortenerStorage, error) {
	switch cfg.Shortener.FileStoragePath {
	case "":
		return NewURLStorage(NewMapStorage()), nil
	default:
		fileStorage, err := NewURLFileStorage(cfg.Shortener.FileStoragePath)
		if err != nil {
			return nil, err
		}
		return NewURLStorage(fileStorage), nil
	}
}
