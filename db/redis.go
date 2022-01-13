package db

import (
	"github.com/blogs/helper"
	"github.com/go-redis/redis/v7"
)

var RedisClient *redis.Client

// InitRedis initializes a redis client
func InitRedis() {
	dsn := "localhost:6379"
	RedisClient = redis.NewClient(&redis.Options{
		Addr: dsn,
		DB:   0,
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
	helper.Logger.Infow("Connected to Redis successfully!")
}
