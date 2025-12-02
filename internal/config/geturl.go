package config

import (
	"os"

	"github.com/joho/godotenv"
)

func GetUrl() string {
	godotenv.Load()

	url := os.Getenv("ORIGIN_URL")

	return url
}