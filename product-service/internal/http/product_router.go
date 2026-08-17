package http

import (
	"net/http"

	"github.com/esmrzv/product-service/internal/handler"
)


func NewProductRouter(productHandler handler.ProductHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /products", productHandler.CreateProduct)
	return mux
}