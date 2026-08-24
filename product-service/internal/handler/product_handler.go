package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/esmrzv/product-service/internal/dto"
	"github.com/esmrzv/product-service/internal/middleware"
	"github.com/esmrzv/product-service/internal/service"
	"github.com/google/uuid"
)

type ProductHandler interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
	GetProductByID(w http.ResponseWriter, r *http.Request)
	GetAllProducts(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
}

type productHandler struct {
	service service.ProductService
}

func NewProductHandler(s service.ProductService) ProductHandler {
	return &productHandler{
		service: s,
	}
}

func (h *productHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "missing user id in context", http.StatusInternalServerError)
		return
	}
	log.Printf("authenticated user: %s", userID)
	response, err := h.service.CreateProduct(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			http.Error(w, "category not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to create product", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

func (h *productHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	id, err := uuid.Parse(productID)
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to get product", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (h *productHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProducts(r.Context())
	if err != nil {
		http.Error(w, "failed to get products", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (h *productHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	id, err := uuid.Parse(productID)
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}
	var req dto.UpdateProductRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateProduct(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (h *productHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	id, err := uuid.Parse(productID)
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "invalid product ID", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
