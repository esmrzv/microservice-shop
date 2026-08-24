package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/esmrzv/product-service/internal/auth"
	"github.com/esmrzv/product-service/internal/config"
	"github.com/esmrzv/product-service/internal/database"
	"github.com/esmrzv/product-service/internal/handler"
	apiHTTP "github.com/esmrzv/product-service/internal/http"
	"github.com/esmrzv/product-service/internal/middleware"
	"github.com/esmrzv/product-service/internal/repository"
	"github.com/esmrzv/product-service/internal/service"
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
		log.Fatal("error connect to db")
	}
	defer db.Close()

	tokenService := auth.NewTokenService(cfg.JWTSecret)
	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	productService := service.NewProductService(productRepo, categoryRepo)
	productHandler := handler.NewProductHandler(productService)

	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	router := apiHTTP.NewRouter(productHandler, categoryHandler, authMiddleware)

	server := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		log.Printf("server starting in port :%s", cfg.HTTPPort)
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

}
