package config

import (
	
	"fmt"
	"os"
)

type Config struct {
	HTTPPort string

	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string

	JWTSecret string
}


func Load() (Config, error){
	cfg := Config{
		HTTPPort: os.Getenv("HTTP_PORT"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
	if err := validateRequired(cfg.DBHost, "DB_HOST"); err != nil{
		return Config{}, err
	}
	if err := validateRequired(cfg.DBPort, "DB_PORT"); err != nil{
		return Config{}, err
	}
	if err := validateRequired(cfg.DBUser, "DB_USER"); err != nil{
		return Config{}, err
	}
	if err := validateRequired(cfg.DBPassword, "DB_PASSWORD"); err != nil{
		return Config{}, err
	}
	if err := validateRequired(cfg.DBName, "DB_NAME"); err != nil{
		return Config{}, err
	}
	if err := validateRequired(cfg.JWTSecret, "JWT_SECRET"); err != nil{
		return Config{}, err
	}
	return cfg, nil

}


func validateRequired(value, field string) error{
	if value == "" {
		return fmt.Errorf("%s is required", field)
	}

	return nil

}