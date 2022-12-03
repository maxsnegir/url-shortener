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
	SaveData(ctx context.Context, userID, url string) (string, error)
	SaveDataBatch(ctx context.Context, userID string, originalURLs []URLDataBatchRequest) ([]URLDataBatchResponse, error)
	GetOriginalURL(ctx context.Context, shortURLID string) (string, error)
	GetUserURLs(ctx context.Context, userID string) ([]storage.URLData, error)
	GetHostURL() string
	IsURLValid(url string) error
	Ping(ctx context.Context) error
}

type shortener struct {
	storage storage.ShortenerStorage
	hostURL string
}

type URLDataBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type URLDataBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (s *shortener) SaveData(ctx context.Context, userID string, url string) (string, error) {
	if err := s.IsURLValid(url); err != nil {
		return "", err
	}
	urlHash := s.getURLHash(url)

	urlData := storage.URLData{
		ShortURL:    fmt.Sprintf("%s/%s/", s.GetHostURL(), urlHash),
		OriginalURL: url,
	}
	if err := s.storage.SaveData(ctx, userID, urlData); err != nil {
		return "", err
	}
	return urlData.ShortURL, nil
}

func (s *shortener) SaveDataBatch(ctx context.Context, userID string, originalURLs []URLDataBatchRequest) ([]URLDataBatchResponse, error) {
	urlDataList := make([]storage.URLData, 0, len(originalURLs))
	urlDataResponse := make([]URLDataBatchResponse, 0, len(originalURLs))

	for _, originalURL := range originalURLs {

		if err := s.IsURLValid(originalURL.OriginalURL); err != nil {
			return urlDataResponse, err
		}
		urlHash := s.getURLHash(originalURL.OriginalURL)

		urlData := storage.URLData{
			ShortURL:    fmt.Sprintf("%s/%s/", s.GetHostURL(), urlHash),
			OriginalURL: originalURL.OriginalURL,
		}
		urlDataList = append(urlDataList, urlData)
		urlDataResponse = append(urlDataResponse, URLDataBatchResponse{
			CorrelationID: originalURL.CorrelationID,
			ShortURL:      urlData.ShortURL,
		})
	}
	err := s.storage.SaveDataBatch(ctx, userID, urlDataList)
	return urlDataResponse, err
}

func (s *shortener) IsURLValid(URL string) error {
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

func (s *shortener) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	originalURL, err := s.storage.GetOriginalURL(ctx, shortURL)
	if err != nil {
		return "", OriginalURLNotFound{shortURL}
	}
	return originalURL, nil
}

func (s *shortener) GetUserURLs(ctx context.Context, userID string) ([]storage.URLData, error) {
	return s.storage.GetUserURLs(ctx, userID)
}

func (s *shortener) Ping(ctx context.Context) error {
	return s.storage.Ping(ctx)
}

func NewShortener(urlStorage storage.ShortenerStorage, hostURL string) URLService {
	return &shortener{
		storage: urlStorage,
		hostURL: hostURL,
	}
}
