package services

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/maxsnegir/url-shortener/internal/storage"
)

type URLService interface {
	SetShortURL(userID, url string) (string, error)
	GetOriginalURLByShort(shortURLID string) (string, error)
	GetAllUserURLs(userID string) ([]storage.URLData, error)
	GetHostURL() string
	Ping(ctx context.Context) error
}

type shortener struct {
	storage storage.ShortenerStorage
	hostURL string
}

func (s *shortener) SetShortURL(userID string, url string) (string, error) {
	if err := s.isURLValid(url); err != nil {
		return "", err
	}
	urlHash := s.getURLHash(url)

	urlData := storage.URLData{
		ShortURL:    fmt.Sprintf("%s/%s/", s.hostURL, urlHash),
		OriginalURL: url,
	}
	if err := s.saveURL(userID, urlData); err != nil {
		return "", err
	}
	return urlData.ShortURL, nil
}

func (s *shortener) isURLValid(URL string) error {
	u, err := url.Parse(URL)
	if err != nil || (u.Scheme == "" || u.Host == "") {
		return URLIsNotValidError{URL: URL}
	}
	return nil
}

func (s *shortener) GetHostURL() string {
	return s.hostURL
}

func (s *shortener) getURLHash(URL string) string {
	hasher := sha1.New()
	hasher.Write([]byte(URL))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha[:8]
}

func (s *shortener) saveURL(userID string, urlData storage.URLData) error {
	if err := s.storage.SetShortURL(urlData); err != nil {
		return err
	}
	if err := s.storage.SetUserURL(userID, urlData.ShortURL); err != nil {
		return err
	}
	return nil
}

func (s *shortener) GetOriginalURLByShort(shortURL string) (string, error) {
	originalURL, err := s.storage.GetOriginalURL(shortURL)
	if err != nil {
		return "", OriginalURLNotFound{shortURL}
	}
	return originalURL, nil
}

func (s *shortener) GetAllUserURLs(userID string) ([]storage.URLData, error) {
	var userURLs []storage.URLData
	shortURLs, err := s.storage.GetUserURLs(userID)
	if err != nil {
		return userURLs, err
	}
	for _, shortURL := range shortURLs {
		originalURL, err := s.storage.GetOriginalURL(shortURL)
		if err != nil {
			continue
		}
		userURLs = append(userURLs, storage.URLData{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		})
	}
	return userURLs, nil
}

func (s *shortener) Ping(ctx context.Context) error {
	return s.storage.Ping(ctx)
}

func NewShortener(urlStorage storage.ShortenerStorage, hostURL string) *shortener {
	return &shortener{
		storage: urlStorage,
		hostURL: hostURL,
	}
}
