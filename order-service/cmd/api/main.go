package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/esmrzv/microservice-shop/order-service/internal/auth"
	"github.com/esmrzv/microservice-shop/order-service/internal/config"
	"github.com/esmrzv/microservice-shop/order-service/internal/database"
	grpcHandler "github.com/esmrzv/microservice-shop/order-service/internal/grpc"
	"github.com/esmrzv/microservice-shop/order-service/internal/handler"
	"github.com/esmrzv/microservice-shop/order-service/internal/middleware"
	"github.com/esmrzv/microservice-shop/order-service/internal/repository"
	"github.com/esmrzv/microservice-shop/order-service/internal/service"
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

	grpcConn, err := grpc.NewClient(cfg.ProductGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer grpcConn.Close()
	productGrpcClient := product.NewProductServiceClient(grpcConn)

	tokenService := auth.NewTokenService(cfg.JwtSECRET)
	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	orderRepo := repository.NewOrderRepository(db)
	productClient := grpcHandler.NewProductClient(productGrpcClient)
	orderService := service.NewOrderService(orderRepo, productClient)
	orderHandler := handler.NewOrderHandler(orderService)
	mux := http.NewServeMux()
	handler.OrderRoutes(mux, orderHandler, authMiddleware)
	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: mux,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("order - server listening on port %s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Printf("context canceled, shutting down server")
	case err := <-serverErr:
		log.Printf("server error: %s", err)
		return
	}
	log.Println("Shutting down signal received")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %s", err)
	}

	db.Close()
	log.Println("order server stopped")

}
