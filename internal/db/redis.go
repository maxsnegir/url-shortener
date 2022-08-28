package db

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/maxsnegir/url-shortener/cmd/config"
)

func NewRedis(cfg config.Config) (*redis.Client, error) {
	redisAdd := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAdd,
		DB:   cfg.Redis.Db,
	})

	if err := rdb.Ping().Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
