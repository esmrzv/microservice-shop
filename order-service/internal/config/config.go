package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort        string
	DBName          string
	DBUser          string
	DBPass          string
	DBHost          string
	DBPort          string
	ProductGRPCAddr string

	JwtSECRET string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &Config{
		HTTPPort:        os.Getenv("HTTP_PORT"),
		DBName:          os.Getenv("DB_NAME"),
		DBUser:          os.Getenv("DB_USER"),
		DBPass:          os.Getenv("DB_PASSWORD"),
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		JwtSECRET:       os.Getenv("JWT_SECRET"),
		ProductGRPCAddr: os.Getenv("PRODUCT_GRPC_ADDR"),
	}
	return config, nil

}
