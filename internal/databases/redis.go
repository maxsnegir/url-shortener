package databases

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/maxsnegir/url-shortener/cmd/config"
)

// NewRedis Изначально я хотел хранить урлы в редис, реализовал для этого логику, но не подумал про
// автотесты, поэтому пусть это лежит до лучших времен :(
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
