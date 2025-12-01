package config

import (
	"os"

	"github.com/joho/godotenv"
)

func GetDatabase() string {
	godotenv.Load()

	dbUrl := os.Getenv("DATABASE_URL")

	return dbUrl
}