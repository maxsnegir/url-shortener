package services

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/maxsnegir/url-shortener/internal/logging"
	"github.com/maxsnegir/url-shortener/internal/storage"
)

type ShortenerService interface {
	SaveData(ctx context.Context, userToken, originalURL string) (string, error)
	SaveDataBatch(ctx context.Context, userToken string, originalURLs []URLDataBatchRequest) ([]URLDataBatchResponse, error)
	GetOriginalURL(ctx context.Context, shortURLID string) (storage.URLData, error)
	GetUserURLs(ctx context.Context, userToken string) ([]storage.URLData, error)
	GetHostURL() string
	IsURLValid(url string) error
	Ping(ctx context.Context) error
	DeleteURLs(urlIDsToDel []string)
	GetURLIdFromShortURL(shortURL string) string
	Shutdown()
}

const BatchSizeForDelete = 100

type shortener struct {
	storage            storage.ShortenerStorage
	hostURL            string
	logger             *logrus.Logger
	urlsToDeleteQueue  chan []string
	done               chan struct{}
	deleteWorkersCount int
}

type URLDataBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type URLDataBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (s *shortener) SaveData(ctx context.Context, userToken string, url string) (string, error) {
	if err := s.IsURLValid(url); err != nil {
		return "", err
	}
	urlHash := s.getURLHash(url)

	urlData := storage.URLData{
		ShortURL:    fmt.Sprintf("%s/%s/", s.GetHostURL(), urlHash),
		OriginalURL: url,
	}
	if err := s.storage.SaveData(ctx, userToken, urlData); err != nil {
		return "", err
	}
	return urlData.ShortURL, nil
}

func (s *shortener) SaveDataBatch(ctx context.Context, userToken string, originalURLs []URLDataBatchRequest) ([]URLDataBatchResponse, error) {
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
	err := s.storage.SaveDataBatch(ctx, userToken, urlDataList)
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

func (s *shortener) GetOriginalURL(ctx context.Context, shortURL string) (storage.URLData, error) {
	urlData, err := s.storage.GetOriginalURL(ctx, shortURL)
	if err != nil {
		return urlData, err
	}
	return urlData, nil
}

func (s *shortener) GetUserURLs(ctx context.Context, userToken string) ([]storage.URLData, error) {
	return s.storage.GetUserURLs(ctx, userToken)
}

func (s *shortener) Ping(ctx context.Context) error {
	return s.storage.Ping(ctx)
}

func (s *shortener) GetURLIdFromShortURL(shortURL string) string {
	u, err := url.Parse(shortURL)
	if err != nil {
		return ""
	}
	return strings.Trim(u.Path, "/")
}

func (s *shortener) DeleteURLs(urlIDsToDel []string) {
	urlsToDelete := make([]string, 0, len(urlIDsToDel))
	for _, urlID := range urlIDsToDel {
		urlsToDelete = append(urlsToDelete, fmt.Sprintf("%s/%s/", s.GetHostURL(), urlID))
	}
	s.writeURLsToQueue(urlsToDelete)
}

// writeURLsToQueue записывает ссылки для удаления в канал
func (s *shortener) writeURLsToQueue(shortURLs []string) {
	select {
	case <-s.done:
		return
	default:
		s.urlsToDeleteQueue <- shortURLs
	}
}

// startDeleteWorkers Запускаем воркеры для, которые удаляют ссылки
func (s *shortener) startDeleteWorkers() {
	wg := &sync.WaitGroup{}
	for i := 0; i < s.deleteWorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.deleteURLs()
		}()
	}
	wg.Wait()
}

func (s *shortener) deleteURLs() {
	for shortURLs := range s.urlsToDeleteQueue {
		if len(shortURLs) > BatchSizeForDelete {
			// Если размер батча превышен, снова грузим остаток ссылок в канал
			nextShortURLs := shortURLs[BatchSizeForDelete:]
			shortURLs = shortURLs[:BatchSizeForDelete]
			go func() {
				s.writeURLsToQueue(nextShortURLs)
			}()
		}
		go func(urlsToDelete []string) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := s.storage.DeleteURLs(ctx, urlsToDelete); err != nil {
				s.logger.Error(err)
			}
		}(shortURLs)

	}
}

// createURLsChan Создаем канал в который будут поступать массивы со ссылками для удаления
func (s *shortener) createURLsChan() chan []string {
	c := make(chan []string)
	go func() {
		// Закрытие канала по сигналу
		<-s.done
		close(s.urlsToDeleteQueue)
	}()
	return c
}

func (s *shortener) Shutdown() {
	close(s.done)
}

func NewShortener(urlStorage storage.ShortenerStorage, hostURL string) ShortenerService {
	s := &shortener{
		storage:            urlStorage,
		hostURL:            hostURL,
		logger:             logging.NewLogger("info"),
		done:               make(chan struct{}),
		deleteWorkersCount: 3,
	}
	s.urlsToDeleteQueue = s.createURLsChan()
	go s.startDeleteWorkers()
	return s
}
