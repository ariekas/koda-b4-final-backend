package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func Redis()  *redis.Client {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("faield to get env", err)
	}

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))

	if err != nil {
		fmt.Println("Failed to parse Redis URL", err)
	}

	fmt.Println(" Redis connected successfully")

    return redis.NewClient(opt)
}