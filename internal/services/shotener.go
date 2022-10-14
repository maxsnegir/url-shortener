package services

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/maxsnegir/url-shortener/internal/storages"
	"net/url"
	"strings"
)

type URLService interface {
	SetURL(url string) (string, error)
	GetURLByID(urlID string) (string, error)
}

type shortener struct {
	storage storages.URLDataBase
	hostURL string
}

func (s *shortener) SetURL(url string) (string, error) {
	if err := s.isURLValid(url); err != nil {
		return "", err
	}
	urlHash := s.getURLHash(url)
	if err := s.saveURL(urlHash, url); err != nil {
		return "", err
	}
	shortURL := fmt.Sprintf("%s/%s/", s.hostURL, urlHash)
	return shortURL, nil
}

func (s *shortener) GetURLByID(urlID string) (string, error) {
	originalURL, err := s.storage.Get(urlID)
	if err == nil {
		return originalURL, nil
	}

	if err == storages.KeyError {
		return "", OriginalURLNotFound{urlID}
	}

	return "", err
}

func (s *shortener) saveURL(urlID, originalURL string) error {
	if err := s.storage.Set(urlID, originalURL); err != nil {
		return err
	}
	return nil
}

func (s *shortener) isURLValid(URL string) error {
	u, err := url.Parse(URL)
	if err != nil || (u.Scheme == "" && u.Host == "") {
		return URLIsNotValidError{URL: u.String()}
	}
	return nil
}

func (s *shortener) getURLHash(URL string) string {
	hasher := sha1.New()
	hasher.Write([]byte(URL))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha[:8]
}

func (s *shortener) getURLIdFromShortURL(shortURL string) string {
	URL, err := url.Parse(shortURL)
	if err != nil {
		return ""
	}
	params := strings.Split(URL.Path, "/")
	if len(params) == 3 {
		return params[1]
	}
	return ""
}

func NewShortener(storage storages.URLDataBase, hostURL string) *shortener {
	return &shortener{
		storage: storage,
		hostURL: hostURL,
	}
}
