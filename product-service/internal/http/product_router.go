package http

import (
	"net/http"

	"github.com/esmrzv/product-service/internal/handler"
	"github.com/esmrzv/product-service/internal/middleware"
)

func NewProductRouter(productHandler handler.ProductHandler,
	authMiddleware *middleware.AuthMiddleware,
) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /products/{id}", productHandler.GetProductByID)
	mux.HandleFunc("GET /products", productHandler.GetAllProducts)
	mux.Handle("POST /products",
		authMiddleware.RequireAuth(
			http.HandlerFunc(
				productHandler.CreateProduct)))

	mux.Handle("PUT /products/{id}",
		authMiddleware.RequireAuth(
			http.HandlerFunc(
				productHandler.UpdateProduct)))

	mux.Handle("DELETE /products/{id}",
		authMiddleware.RequireAuth(
			http.HandlerFunc(
				productHandler.DeleteProduct)))
	return mux
}
