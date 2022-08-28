package services

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/go-redis/redis"
	"net/url"
	"time"
)

type UrlIsNotValidError struct {
	Url string
}

func (e UrlIsNotValidError) Error() string {
	return fmt.Sprintf("Url %s is not valid", e.Url)
}

type Shortener struct {
	redisClient *redis.Client
}

func (s *Shortener) SetUrl(originalUrl string, expires time.Duration) (string, error) {
	urlHash := s.hashUrl(originalUrl)
	if err := s.setUrl(urlHash, originalUrl, expires); err != nil {
		return "", err
	}
	return urlHash, nil
}

func (s *Shortener) GetUrlById(urlId string) (string, error) {
	originalUrl, err := s.redisClient.Get(urlId).Result()
	if err != nil {
		return "", err
	}
	return originalUrl, nil
}

func (s *Shortener) setUrl(urlHash, url string, expires time.Duration) error {
	status := s.redisClient.Set(urlHash, url, expires*time.Minute)
	if err := status.Err(); err != nil {
		return err
	}
	return nil
}

func (s *Shortener) UrlIsValid(Url string) bool {
	u, err := url.Parse(Url)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (s *Shortener) hashUrl(originalUrl string) string {
	hasher := sha1.New()
	hasher.Write([]byte(originalUrl))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha[:8]
}
func NewShortener(redisClient *redis.Client) *Shortener {
	return &Shortener{
		redisClient: redisClient,
	}
}
