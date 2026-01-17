package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDSN string
}

func Load() Config {
	// Load .env file (only for local dev)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
	}
}

