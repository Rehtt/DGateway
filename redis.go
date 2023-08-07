package main

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis() error {
	rdb = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   1,
	})
	return rdb.Ping(context.Background()).Err()
}
