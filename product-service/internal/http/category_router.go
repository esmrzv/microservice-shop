package http

import (
	"net/http"

	"github.com/esmrzv/product-service/internal/handler"
)

func NewCategoryRouter(categoryHandler handler.CategoryHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /categories", categoryHandler.CreateCategory)
	mux.HandleFunc("GET /categories/{id}", categoryHandler.GetCategoryByID)
	return mux
}
