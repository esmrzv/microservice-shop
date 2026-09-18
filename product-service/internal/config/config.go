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

	DBName    string
	JWTSecret string

	RedisHost    string
	RedisPort    string
	RedisEnabled bool

	GRPCPort string
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		HTTPPort:     os.Getenv("HTTP_PORT"),
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBName:       os.Getenv("DB_NAME"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		RedisHost:    os.Getenv("REDIS_HOST"),
		RedisPort:    os.Getenv("REDIS_PORT"),
		RedisEnabled: os.Getenv("REDIS_ENABLED") == "true",
		GRPCPort:     os.Getenv("GRPC_PORT"),
	}
	return cfg, nil
}
