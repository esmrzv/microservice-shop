package handler

import (
	"encoding/json"
	"net/http"

	"github.com/esmrzv/product-service/internal/dto"
	"github.com/esmrzv/product-service/internal/service"
)


type ProductHandler interface {
	CreateProduct(w http.ResponseWriter, r *http.Request) 
}


type productHandler struct {
	service service.ProductService
}


func NewProductHandler(s service.ProductService) ProductHandler {
	return &productHandler{
		service: s,
	}
}

func (h *productHandler) CreateProduct(w http.ResponseWriter, r *http.Request){
	var req dto.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil{
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateProduct(r.Context(), req)
	if err != nil{
		http.Error(w, "failed to create product", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	

}