package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/esmrzv/product-service/internal/dto"
	"github.com/esmrzv/product-service/internal/service"
	"github.com/google/uuid"
)

type CategoryHandler interface {
	CreateCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryByID(w http.ResponseWriter, r *http.Request)
}

type categoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) CategoryHandler {
	return &categoryHandler{
		service: service,
	}
}

// @Summary Create category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body dto.CategoryRequest true "Category data"
// @Success 201 {object} models.Category
// @Failure 400 {string} string "Invalid category"
// @Failure 500 {string} string "Internal server error"
// @Router /categories [post]
func (h *categoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateCategory handler called")
	category := dto.CategoryRequest{}
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusInternalServerError)
		return
	}
	result, err := h.service.CreateCategory(r.Context(), category.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// @Summary Get category by ID
// @Description Get a category by its UUID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} models.Category
// @Failure 400 {string} string "Invalid category ID"
// @Failure 404 {string} string "Category not found"
// @Failure 500 {string} string "Internal server error"
// @Router /categories/{id} [get]
func (h *categoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	categoryID := r.PathValue("id")
	id, err := uuid.Parse(categoryID)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid category id: %s", categoryID), http.StatusBadRequest)
		return
	}
	category, err := h.service.GetCategoryByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			http.Error(w, "category not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to retrieve category", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)

}
