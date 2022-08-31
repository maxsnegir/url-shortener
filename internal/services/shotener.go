package services

import (
	"crypto/sha1"
	"encoding/base64"
	"github.com/maxsnegir/url-shortener/internal/storages"
	"net/url"
	"time"
)

type Shortener struct {
	db storages.KeyValueStorage
}

func (s *Shortener) SetURL(originalURL *url.URL, expires time.Duration) (string, error) {
	urlHash := s.hashURL(originalURL.String())
	if err := s.setURL(urlHash, originalURL.String(), expires); err != nil {
		return "", err
	}
	return urlHash, nil
}

func (s *Shortener) GetURLByID(urlID string) (string, error) {
	original, err := s.db.Get(urlID)
	if err != nil {
		if err == storages.KeyError {
			return "", OriginalURLNotFound{urlID}
		}
		return "", err
	}
	originalURL, _ := original.(string)
	return originalURL, nil
}

func (s *Shortener) setURL(urlHash, url string, expires time.Duration) error {
	err := s.db.Set(urlHash, url, expires*time.Minute)
	if err != nil {
		return err
	}
	return nil
}

func (s *Shortener) ParseURL(URL string) (*url.URL, error) {
	u, err := url.Parse(URL)
	if err != nil || (u.Scheme == "" && u.Host == "") {
		return nil, URLIsNotValidError{URL: u.String()}
	}
	return u, nil
}

func (s *Shortener) hashURL(URL string) string {
	hasher := sha1.New()
	hasher.Write([]byte(URL))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha[:8]
}
func NewShortener(db storages.KeyValueStorage) *Shortener {
	return &Shortener{
		db: db,
	}
}
