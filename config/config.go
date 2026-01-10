package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env not loaded, using system env")
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
