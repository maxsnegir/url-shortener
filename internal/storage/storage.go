package storage

import "github.com/maxsnegir/url-shortener/cmd/config"

type KeyValue struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

type URLData struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLDataList []URLData

type Storage interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Shutdown() error
}

type URLStorage interface {
	Storage
	SetURLData(userID string, data URLData) error
	GetURLData(userID string) (URLDataList, error)
}

func GetURLStorage(cfg config.Config) (URLStorage, error) {
	mapStorage := NewURLMapStorage()

	switch cfg.Shortener.FileStoragePath {
	case "":
		return mapStorage, nil
	default:
		fileStorage, err := NewURLFileStorage(cfg.Shortener.FileStoragePath, mapStorage)
		if err != nil {
			return fileStorage, err
		}
		return fileStorage, nil
	}
}
