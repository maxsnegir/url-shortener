package storage

import (
	"context"
	"encoding/json"
)

type URLStorage struct {
	userURLStorage Storage
	urlStorage     Storage
}

func (s *URLStorage) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	encodedData, err := s.urlStorage.Get(shortURL)
	if err != nil {
		return "", nil
	}
	return string(encodedData), nil
}

func (s *URLStorage) SetShortURL(urlData URLData) error {
	if _, err := s.urlStorage.Get(urlData.ShortURL); err == nil {
		// Имитация ошибки при существующей ссылки в базе
		return NewDuplicateError(urlData.ShortURL)
	}
	return s.urlStorage.Set(urlData.ShortURL, []byte(urlData.OriginalURL))
}

func (s *URLStorage) getUserShortURLs(userToken string) ([]string, error) {
	var shortURLs []string

	encodedURLs, err := s.userURLStorage.Get(userToken)
	if err != nil {
		return shortURLs, nil
	}
	err = json.Unmarshal(encodedURLs, &shortURLs)
	return shortURLs, err
}

func (s *URLStorage) GetUserURLs(ctx context.Context, userToken string) ([]URLData, error) {
	var userURLData []URLData

	shortURLs, err := s.getUserShortURLs(userToken)
	if err != nil {
		return userURLData, err
	}

	for _, shortURL := range shortURLs {
		originalURL, err := s.GetOriginalURL(context.TODO(), shortURL) //TODO
		if err != nil {
			continue
		}
		userURLData = append(userURLData, URLData{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		})
	}
	return userURLData, nil
}

func (s *URLStorage) SetUserURL(userToken string, shortURL string) error {
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

func (s *URLStorage) SaveData(ctx context.Context, userToken string, urlData URLData) error {
	if err := s.SetShortURL(urlData); err != nil {
		return err
	}
	if err := s.SetUserURL(userToken, urlData.ShortURL); err != nil {
		return err
	}
	return nil
}

func (s *URLStorage) SaveDataBatch(ctx context.Context, userToken string, urlData []URLData) (err error) {
	for _, url := range urlData {
		if err := s.SaveData(ctx, userToken, url); err != nil {
			return err
		}
	}
	return nil
}

func (s *URLStorage) Ping(ctx context.Context) error {
	return nil
}

func (s *URLStorage) Shutdown(ctx context.Context) error {
	if err := s.urlStorage.Shutdown(ctx); err != nil {
		return err
	}
	return s.userURLStorage.Shutdown(ctx)
}

func NewURLStorage(urlStorage Storage) *URLStorage {
	return &URLStorage{
		urlStorage:     urlStorage,
		userURLStorage: NewMapStorage(), // InMemoryStorage по-дефолту
	}
}
