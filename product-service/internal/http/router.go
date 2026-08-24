package http

import (
	"net/http"

	"github.com/esmrzv/product-service/internal/handler"
	"github.com/esmrzv/product-service/internal/middleware"
)

func NewRouter(
	productHandler handler.ProductHandler,
	categoryHandler handler.CategoryHandler,
	authMiddleware *middleware.AuthMiddleware,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("POST /products",
		authMiddleware.RequireAuth(
			http.HandlerFunc(productHandler.CreateProduct),
		),
	)

	mux.Handle("GET /products/{id}",
		authMiddleware.RequireAuth(
			http.HandlerFunc(productHandler.GetProductByID),
		),
	)

	mux.Handle("GET /products",
		authMiddleware.RequireAuth(
			http.HandlerFunc(productHandler.GetAllProducts),
		),
	)

	mux.Handle("PUT /products/{id}",
		authMiddleware.RequireAuth(
			http.HandlerFunc(productHandler.UpdateProduct),
		),
	)

	mux.Handle("DELETE /products/{id}",
		authMiddleware.RequireAuth(
			http.HandlerFunc(productHandler.DeleteProduct),
		),
	)

	mux.HandleFunc("POST /categories", categoryHandler.CreateCategory)
	mux.HandleFunc("GET /categories/{id}", categoryHandler.GetCategoryByID)

	return mux
}
