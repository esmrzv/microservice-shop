package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort string

	DBHost string
	DBPort string

	DBUser     string
	DBPassword string

	DBName string
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		HTTPPort:   os.Getenv("HTTP_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}
	return cfg, nil
}
