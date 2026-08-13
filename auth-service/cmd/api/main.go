package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/esmrzv/microservice-shop/auth-service/internal/auth"
	"github.com/esmrzv/microservice-shop/auth-service/internal/config"
	"github.com/esmrzv/microservice-shop/auth-service/internal/database"
	"github.com/esmrzv/microservice-shop/auth-service/internal/handler"
	apphttp "github.com/esmrzv/microservice-shop/auth-service/internal/http"
	"github.com/esmrzv/microservice-shop/auth-service/internal/repository"
	"github.com/esmrzv/microservice-shop/auth-service/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

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
	router := apphttp.NewRouter(authHandler, authMiddleware)

	server := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: router,
	}
	go func() {
		log.Printf("Listening on port %s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down http server...")
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctxShutDown); err != nil {
		log.Printf("shutdown: %s\n", err)
	}
}
