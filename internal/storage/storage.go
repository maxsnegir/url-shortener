package storage

import "github.com/maxsnegir/url-shortener/cmd/config"

type Storage interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Shutdown() error
}

func GetStorage(cfg config.Config) (Storage, error) {
	switch cfg.Shortener.FileStoragePath {
	case "":
		return NewMapURLDataBase(), nil
	default:
		return NewFileStorage(cfg.Shortener.FileStoragePath)
	}
}
