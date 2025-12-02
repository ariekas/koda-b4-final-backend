package config

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)
var RedisClient *redis.Client

func InitRedis() {
	godotenv.Load()

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic("Failed to parse Redis URL: " + err.Error())
	}

	RedisClient = redis.NewClient(opt)

	_, err = RedisClient.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	fmt.Println("Redis connected successfully")
}