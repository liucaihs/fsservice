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
		PoolSize:   10,
		MaxRetries: 3,
	})

	return rdsclt.Ping().Err()
}

func DatabaseClose() {
	if err := rdsclt.Close(); err != nil {
		LogErr("Err from storage.dbinit.DatabaseClose(): ", err)
	}
}
