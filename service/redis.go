package service

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func connectToRedis(c *Config) (*redis.Client, error) {
	if c == nil {
		return nil, ErrInternalError
	}
	if c.Store.Redis == nil {
		return nil, ErrInvalidConfig
	}
	if c.Store.Redis.Host == "" {
		logrus.Warn("Redis host is empty, falling back to default")
		c.Store.Redis.Host = "127.0.0.1:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr:     c.Store.Redis.Host,
		Password: c.Store.Redis.Password,
		DB:       c.Store.Redis.DB,
	})
	return client, client.Ping().Err()
}
