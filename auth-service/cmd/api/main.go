package main

import (
	"context"
	"log"
	"net/http"

	"github.com/esmrzv/microservice-shop/auth-service/internal/auth"
	"github.com/esmrzv/microservice-shop/auth-service/internal/config"
	"github.com/esmrzv/microservice-shop/auth-service/internal/database"
	"github.com/esmrzv/microservice-shop/auth-service/internal/handler"
	"github.com/esmrzv/microservice-shop/auth-service/internal/repository"
	"github.com/esmrzv/microservice-shop/auth-service/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	db, err := database.NewPostgres(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	authMiddleware := auth.NewAuthMiddleware(cfg.JWTSecret)
	tokenService := auth.NewTokenService(cfg.JWTSecret)
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, tokenService)
	authHandler := handler.NewAuthHandler(authService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.Handle("/me", authMiddleware.RequireAuth(http.HandlerFunc(authHandler.Me)))
	server := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: mux,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on port", cfg.HTTPPort)

}
