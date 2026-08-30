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

// @Summary Create product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.CreateProductRequest true "Product data"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Category not found"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /products [post]
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
	response, err := h.service.CreateProduct(r.Context(), userID, req)
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

// @Summary Get product by ID
// @Description Get a product by its UUID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {string} string "Invalid product ID"
// @Failure 404 {string} string "Product not found"
// @Failure 500 {string} string "failed to get product"
// @Security BearerAuth
// @Router /products/{id} [get]
func (h *productHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	id, err := uuid.Parse(productID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get product", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// @Summary Get all products
// @Description Get list of all products
// @Tags products
// @Produce json
// @Success 200 {array} dto.ProductResponse
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /products [get]
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

// @Summary Update product
// @Description Update an existing product by UUID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body dto.UpdateProductRequest true "Updated product data"
// @Success 204
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Product not found"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /products/{id} [put]
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
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "missing user id in context", http.StatusInternalServerError)
		return
	}
	err = h.service.UpdateProduct(r.Context(), userID, id, &req)
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

// @Summary Delete product
// @Description Delete a product by UUID
// @Tags products
// @Param id path string true "Product ID"
// @Success 204
// @Failure 400 {string} string "Invalid product ID"
// @Failure 404 {string} string "Product not found"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /products/{id} [delete]
func (h *productHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	id, err := uuid.Parse(productID)
	if err != nil {
		http.Error(w, "failed to parse uuid", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "missing user id in context", http.StatusInternalServerError)
		return
	}
	err = h.service.DeleteProduct(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
