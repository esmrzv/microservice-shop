package http

import (
	"net/http"

	"github.com/esmrzv/product-service/internal/handler"
)

func NewProductRouter(productHandler handler.ProductHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /products", productHandler.CreateProduct)
	mux.HandleFunc("GET /products/{id}", productHandler.GetProductByID)
	mux.HandleFunc("GET /products", productHandler.GetAllProducts)
	mux.HandleFunc("PUT /products/{id}", productHandler.UpdateProduct)
	mux.HandleFunc("DELETE /products/{id}", productHandler.DeleteProduct)
	return mux
}
