package handler

import (
	"net/http"

	"github.com/esmrzv/product-service/internal/service"
)

func setupHandlerRoute(service service.ProductService) *http.ServeMux {
	handler := NewProductHandler(service)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", handler.GetAllProducts)
	mux.HandleFunc("GET /products/{id}", handler.GetProductByID)
	mux.HandleFunc("POST /products", handler.CreateProduct)
	mux.HandleFunc("PUT /products/{id}", handler.UpdateProduct)
	mux.HandleFunc("DELETE /products/{id}", handler.DeleteProduct)

	return mux

}
