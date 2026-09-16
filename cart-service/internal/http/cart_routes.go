package http

import (
	"net/http"

	"github.com/esmrzv/cart-service/internal/handler"
	"github.com/esmrzv/cart-service/internal/middleware"
)

func CartRoutes(
	mux *http.ServeMux,
	cartHandler handler.CartHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	mux.Handle(
		"GET /cart", authMiddleware.RequireAuth(
			http.HandlerFunc(cartHandler.GetCart)))
	mux.Handle("POST /cart/items", authMiddleware.RequireAuth(
		http.HandlerFunc(cartHandler.AddItem)))
	mux.Handle("PUT /cart/items/{id}", authMiddleware.RequireAuth(
		http.HandlerFunc(cartHandler.UpdateItemQuantity)))
	mux.Handle("DELETE /cart/items/{id}", authMiddleware.RequireAuth(
		http.HandlerFunc(cartHandler.DeleteItem)))
	mux.Handle("DELETE /cart", authMiddleware.RequireAuth(
		http.HandlerFunc(cartHandler.DeleteCart)))

}
