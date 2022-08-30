package services

import (
	"crypto/sha1"
	"encoding/base64"
	"github.com/maxsnegir/url-shortener/internal/databases"
	"net/url"
	"time"
)

type Shortener struct {
	db databases.KeyValueDB
}

func (s *Shortener) SetURL(originalURL string, expires time.Duration) (string, error) {
	urlHash := s.hashUrl(originalURL)
	if err := s.setURL(urlHash, originalURL, expires); err != nil {
		return "", err
	}
	return urlHash, nil
}

func (s *Shortener) GetURLByID(urlID string) (string, error) {
	original, err := s.db.Get(urlID)
	if err != nil {
		if err == databases.KeyError {
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

func (s *Shortener) URLIsValid(Url string) bool {
	u, err := url.Parse(Url)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (s *Shortener) hashUrl(originalUrl string) string {
	hasher := sha1.New()
	hasher.Write([]byte(originalUrl))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha[:8]
}
func NewShortener(db databases.KeyValueDB) *Shortener {
	return &Shortener{
		db: db,
	}
}
