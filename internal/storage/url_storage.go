package storage

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/sync/errgroup"
)

// MemoryURLStorage обертка над MapStorage/FileStorage.
// Используется для избежания дублирования кода, тк и то и другое хранилище имеют общее поведение сохранения данных,
// но формат и структура хранения данных отличается
type MemoryURLStorage struct {
	userURLStorage Storage
	urlStorage     Storage
}

func (s *MemoryURLStorage) GetOriginalURL(ctx context.Context, shortURL string) (URLData, error) {
	urlData := URLData{
		ShortURL: shortURL,
	}
	encodedData, err := s.urlStorage.Get(shortURL)
	if err != nil {
		if errors.Is(err, KeyError) {
			return urlData, NewOriginalURLNotFound(shortURL)
		}
		return urlData, err
	}
	if err := json.Unmarshal(encodedData, &urlData); err != nil {
		return urlData, err
	}
	return urlData, nil
}

func (s *MemoryURLStorage) SetShortURL(urlData URLData) error {
	if _, err := s.urlStorage.Get(urlData.ShortURL); err == nil {
		// Имитация ошибки при существующей ссылки в базе
		return NewDuplicateError(urlData.ShortURL)
	}
	return s.SaveShortURL(urlData)
}

func (s *MemoryURLStorage) SaveShortURL(urlData URLData) error {
	b, err := json.Marshal(urlData)
	if err != nil {
		return err
	}
	return s.urlStorage.Set(urlData.ShortURL, b)
}

func (s *MemoryURLStorage) getUserShortURLs(userToken string) ([]string, error) {
	var shortURLs []string

	encodedURLs, err := s.userURLStorage.Get(userToken)
	if err != nil {
		return shortURLs, nil
	}
	err = json.Unmarshal(encodedURLs, &shortURLs)
	return shortURLs, err
}

func (s *MemoryURLStorage) GetUserURLs(ctx context.Context, userToken string) ([]URLData, error) {
	var userURLData []URLData

	shortURLs, err := s.getUserShortURLs(userToken)
	if err != nil {
		return userURLData, err
	}

	for _, shortURL := range shortURLs {
		urlData, err := s.GetOriginalURL(context.TODO(), shortURL) //TODO
		if err != nil {
			continue
		}
		userURLData = append(userURLData, urlData)
	}
	return userURLData, nil
}

func (s *MemoryURLStorage) SetUserURL(userToken string, shortURL string) error {
	userShortURLs, err := s.getUserShortURLs(userToken)
	if err != nil {
		return err
	}
	for _, url := range userShortURLs {
		if url == shortURL {
			return nil
		}
	}
	userShortURLs = append(userShortURLs, shortURL)
	encodedURLs, err := json.Marshal(userShortURLs)
	if err != nil {
		return err
	}
	return s.userURLStorage.Set(userToken, encodedURLs)
}

func (s *MemoryURLStorage) SaveData(ctx context.Context, userToken string, urlData URLData) error {
	if err := s.SetShortURL(urlData); err != nil {
		return err
	}
	if err := s.SetUserURL(userToken, urlData.ShortURL); err != nil {
		return err
	}
	return nil
}

func (s *MemoryURLStorage) SaveDataBatch(ctx context.Context, userToken string, urlData []URLData) (err error) {
	for _, url := range urlData {
		if err := s.SaveData(ctx, userToken, url); err != nil {
			return err
		}
	}
	return nil
}

func (s *MemoryURLStorage) Ping(ctx context.Context) error {
	return nil
}

func (s *MemoryURLStorage) Shutdown() error {
	if err := s.urlStorage.Shutdown(); err != nil {
		return err
	}
	return s.userURLStorage.Shutdown()
}

func (s *MemoryURLStorage) DeleteURLs(ctx context.Context, urlsToDelete []string) error {
	g, ctx := errgroup.WithContext(context.Background())
	for _, shortURL := range urlsToDelete {
		shortURL := shortURL
		g.Go(func() error {
			urlData, err := s.GetOriginalURL(ctx, shortURL)
			if err != nil {
				return err
			}
			urlData.Deleted = true
			return s.SaveShortURL(urlData)
		})
	}
	return g.Wait()
}

func NewMemoryURLStorage(urlStorage Storage) *MemoryURLStorage {
	return &MemoryURLStorage{
		urlStorage:     urlStorage,
		userURLStorage: NewMapStorage(), // InMemoryStorage по-дефолту
	}
}
