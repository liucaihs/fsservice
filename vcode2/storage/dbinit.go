package storage

import (
	"github.com/go-redis/redis"
)

var rdsclt *redis.Client

func DatabaseInit() error {
	rdsclt = redis.NewClient(&redis.Options{
		Addr:       "redis_server:6379",
		Password:   "",
		DB:         0,
		PoolSize:   999,
		MaxRetries: 3,
	})
	err := rdsclt.Ping().Err()
	return err
}
