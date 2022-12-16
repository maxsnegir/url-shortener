package storage

import (
	"context"
	"sync"
)

// MapStorage In-Memory хранилище
type MapStorage struct {
	storage map[string][]byte
	mutex   *sync.RWMutex
}

func (s *MapStorage) Set(key string, value []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.storage[key] = value
	return nil
}

func (s *MapStorage) Get(key string) ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	value, ok := s.storage[key]
	if !ok {
		return nil, KeyError
	}
	return value, nil
}

func (s *MapStorage) Shutdown(ctx context.Context) error {
	return nil
}

func configureMapStorage() *MapStorage {
	return &MapStorage{storage: map[string][]byte{}, mutex: &sync.RWMutex{}}
}

func NewMapStorage() Storage {
	return configureMapStorage()
}
