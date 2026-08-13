package http

import (
	"net/http"

	"github.com/esmrzv/microservice-shop/auth-service/internal/auth"
	"github.com/esmrzv/microservice-shop/auth-service/internal/handler"
)

func NewRouter(authHandler handler.AuthHandler,
	authMiddleware auth.AuthMiddleware) http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.Handle("/me", authMiddleware.RequireAuth(http.HandlerFunc(authHandler.Me)))
	return mux
}
