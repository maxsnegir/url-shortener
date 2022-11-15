package services

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/maxsnegir/url-shortener/internal/storage"
	"net/url"
	"strings"
)

type URLService interface {
	SetShortURL(userID, url string) (string, error)
	GetOriginalURLByShort(urlID string, shortURLID string) (string, error)
	GetAllUserURLs(userID string) ([]storage.URLData, error)
	GetHostURL() string
}

type shortener struct {
	storage storage.URLStorage
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

func (s *shortener) saveURL(urlID string, urlData storage.URLData) error {
	if err := s.storage.SetURLData(urlID, urlData); err != nil {
		return err
	}
	return nil
}

func (s *shortener) GetOriginalURLByShort(userID string, shortURL string) (string, error) {
	userURls, err := s.storage.GetURLData(userID)
	if err != nil {
		return "", err
	}
	for _, urlData := range userURls {
		if urlData.ShortURL == shortURL {
			return urlData.OriginalURL, nil
		}
	}
	return "", OriginalURLNotFound{shortURL}
}

func (s *shortener) GetAllUserURLs(userID string) ([]storage.URLData, error) {
	return s.storage.GetURLData(userID)
}

func (s *shortener) getURLIdFromShortURL(shortURL string) (urlID string) {
	urlID = "/"
	URL, err := url.Parse(shortURL)
	if err != nil {
		return urlID
	}
	params := strings.Split(URL.Path, "/")
	if len(params) == 3 {
		return urlID + params[1]
	}
	return urlID
}

func NewShortener(storage storage.URLStorage, hostURL string) *shortener {
	return &shortener{
		storage: storage,
		hostURL: hostURL,
	}
}
