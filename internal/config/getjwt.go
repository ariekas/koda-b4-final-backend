package config

import (
	"os"

	"github.com/joho/godotenv"
)

func GetJwtToken() string {
	godotenv.Load()
	JWTtoken  := os.Getenv("JWT_TOKEN")

	return JWTtoken
}