package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort string
	DBHost   string
	DBPort   string

	DBUser     string
	DBPassword string

	DBName          string
	JWTSecret       string
	ProductGRPCAddr string
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		DBName:          os.Getenv("DB_NAME"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		HTTPPort:        os.Getenv("HTTP_PORT"),
		ProductGRPCAddr: os.Getenv("PRODUCT_GRPC_ADDR"),
	}
	return cfg, nil
}
