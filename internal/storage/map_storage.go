package storage

import "context"

// MapStorage In-Memory хранилище
type MapStorage map[string][]byte

func (db MapStorage) Set(key string, value []byte) error {
	db[key] = value
	return nil
}

func (db MapStorage) Get(key string) ([]byte, error) {
	value, ok := db[key]
	if !ok {
		return nil, KeyError
	}
	return value, nil
}

func (db MapStorage) Shutdown(ctx context.Context) error {
	return nil
}

func configureMapStorage() MapStorage {
	return make(MapStorage)
}

func NewMapStorage() Storage {
	return configureMapStorage()
}
