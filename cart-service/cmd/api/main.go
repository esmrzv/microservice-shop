package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/esmrzv/cart-service/internal/config"
	"github.com/esmrzv/cart-service/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := database.NewPostgres(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("database connected")
}
