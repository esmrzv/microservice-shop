package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/esmrzv/product-service/docs"
	"github.com/esmrzv/product-service/internal/auth"
	"github.com/esmrzv/product-service/internal/cache"
	"github.com/esmrzv/product-service/internal/config"
	"github.com/esmrzv/product-service/internal/database"
	grpcHandler "github.com/esmrzv/product-service/internal/grpc"
	"github.com/esmrzv/product-service/internal/handler"
	apiHTTP "github.com/esmrzv/product-service/internal/http"
	"github.com/esmrzv/product-service/internal/middleware"
	"github.com/esmrzv/product-service/internal/repository"
	"github.com/esmrzv/product-service/internal/service"
	"github.com/esmrzv/product-service/proto"
	"google.golang.org/grpc"
)

// @title Product Service API
// @version 1.0
// @description API для управления товарами и категориями.
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
		log.Fatal("error connect to db")
	}
	defer db.Close()

	clientRedis, err := cache.NewRedis(ctx, cfg)
	if err != nil {
		log.Fatal("error connect to redis")
	}
	defer clientRedis.Close()

	tokenService := auth.NewTokenService(cfg.JWTSecret)
	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	productCache := cache.NewProductCache(clientRedis)
	productService := service.NewProductService(productRepo, categoryRepo, productCache)
	productHandler := handler.NewProductHandler(productService)

	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	productServer := grpcHandler.NewProductServer(productService)
	grpcServer := grpc.NewServer()
	proto.RegisterProductServiceServer(grpcServer, productServer)
	grpcListener, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Printf("grpc server starting on %s", cfg.GRPCPort)

		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	router := apiHTTP.NewRouter(productHandler, categoryHandler, authMiddleware)

	server := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		log.Printf("server starting in port %s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutDown); err != nil {
		log.Printf("error shutDown %s\n", err)
	}

	grpcServer.GracefulStop()
	grpcListener.Close()

}
