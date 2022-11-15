package storage

import "encoding/json"

type URLStorage struct {
	userURLStorage Storage
	urlStorage     Storage
}

func (s *URLStorage) GetOriginalURL(shortURL string) (string, error) {
	encodedData, err := s.urlStorage.Get(shortURL)
	if err != nil {
		return "", nil
	}
	return string(encodedData), nil
}

func (s *URLStorage) SetShortURL(urlData URLData) error {
	return s.urlStorage.Set(urlData.ShortURL, []byte(urlData.OriginalURL))
}

func (s *URLStorage) GetUserURLs(userID string) ([]string, error) {
	var userURLs []string
	encodedURLs, err := s.userURLStorage.Get(userID)
	if err != nil {
		return userURLs, nil
	}
	if err := json.Unmarshal(encodedURLs, &userURLs); err != nil {
		return userURLs, err
	}
	return userURLs, nil
}

func (s *URLStorage) SetUserURL(userID string, shortURL string) error {
	userURLs, err := s.GetUserURLs(userID)
	if err != nil {
		return err
	}
	for _, url := range userURLs {
		if url == shortURL {
			return nil
		}
	}
	userURLs = append(userURLs, shortURL)
	encodedURLs, err := json.Marshal(userURLs)
	if err != nil {
		return err
	}
	return s.userURLStorage.Set(userID, encodedURLs)
}

func (s *URLStorage) Shutdown() error {
	if err := s.urlStorage.Shutdown(); err != nil {
		return err
	}
	return s.userURLStorage.Shutdown()
}

func NewURLStorage(urlStorage Storage) *URLStorage {
	return &URLStorage{
		urlStorage:     urlStorage,
		userURLStorage: NewMapStorage(), // InMemoryStorage по-дефолту
	}
}
