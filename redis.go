package main

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis(addr, username, password string, db int) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		Username: username,
		Password: password,
	})
	return rdb.Ping(context.Background()).Err()
}
