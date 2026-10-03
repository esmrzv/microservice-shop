package handler

import (
	"net/http"

	"github.com/esmrzv/microservice-shop/order-service/internal/middleware"
)

func OrderRoutes(mux *http.ServeMux, orderHandler OrderHandler, auth *middleware.AuthMiddleware) {
	mux.Handle("POST /orders", auth.RequireAuth(http.HandlerFunc(orderHandler.Create)))
	mux.Handle("GET /orders/{id}", auth.RequireAuth(http.HandlerFunc(orderHandler.GetByID)))
	mux.Handle("GET /orders", auth.RequireAuth(http.HandlerFunc(orderHandler.GetByUserID)))
}
