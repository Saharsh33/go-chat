package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDSN string
	JWTSecret   string
}

func Load() Config {
	// Load .env file (only for local dev)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required but not set")
	}

	return Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
		JWTSecret:   jwtSecret,
	}
}

