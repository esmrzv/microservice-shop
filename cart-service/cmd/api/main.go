package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/esmrzv/cart-service/internal/auth"
	"github.com/esmrzv/cart-service/internal/config"
	"github.com/esmrzv/cart-service/internal/database"
	grpcHabdler "github.com/esmrzv/cart-service/internal/grpc"
	"github.com/esmrzv/cart-service/internal/handler"
	apihttp "github.com/esmrzv/cart-service/internal/http"
	"github.com/esmrzv/cart-service/internal/middleware"
	"github.com/esmrzv/cart-service/internal/repository"
	"github.com/esmrzv/cart-service/internal/service"
	"github.com/esmrzv/microservice-shop/proto/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	log.Println("database connected")

	grpcConn, err := grpc.NewClient(cfg.ProductGRPCAddr, grpc.WithTransportCredentials(
		insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer grpcConn.Close()

	productGRPCClient := product.NewProductServiceClient(grpcConn)
	productClient := grpcHabdler.NewProductClient(productGRPCClient)

	tokenService := auth.NewTokenService(cfg.JWTSecret)
	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	cartRepo := repository.NewCartRepository(db)
	cartService := service.NewCartService(cartRepo, productClient)
	cartHandler := handler.NewCartHandler(cartService)

	mux := http.NewServeMux()
	apihttp.CartRoutes(mux, cartHandler, authMiddleware)
	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: mux,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("cart-service listening on :%s", cfg.HTTPPort)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {

			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")

	case err := <-serverErr:
		log.Printf("server error: %v", err)
		return
	}

	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	db.Close()
	log.Println("cart-service stopped")

}
