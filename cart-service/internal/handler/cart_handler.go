package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/esmrzv/cart-service/internal/handler/dto"
	"github.com/esmrzv/cart-service/internal/middleware"
	"github.com/esmrzv/cart-service/internal/service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CartHandler interface {
	GetCart(w http.ResponseWriter, r *http.Request)
	AddItem(w http.ResponseWriter, r *http.Request)
	UpdateItemQuantity(w http.ResponseWriter, r *http.Request)
	DeleteItem(w http.ResponseWriter, r *http.Request)
	DeleteCart(w http.ResponseWriter, r *http.Request)
}

type cartHandler struct {
	service service.CartService
}

func NewCartHandler(service service.CartService) CartHandler {
	return &cartHandler{
		service: service,
	}
}

func (h *cartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	cart, err := h.service.GetCart(r.Context(), userID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := dto.CartResponse{
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
		Items:     make([]dto.CartItemResponse, 0, len(cart.Items)),
	}
	for _, item := range cart.Items {
		response.Items = append(response.Items, dto.CartItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

func (h *cartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var req dto.AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body request", http.StatusBadRequest)
		return
	}
	item, err := h.service.AddItem(r.Context(), userID, req.ProductID, req.Quantity)
	if err != nil {
		if errors.Is(err, service.ErrInvalidQuantity) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrProductNotFound) {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := dto.CartItemResponse{
		ID:        item.ID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}

}

func (h *cartHandler) UpdateItemQuantity(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	itemID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body request", http.StatusBadRequest)
		return
	}
	item, err := h.service.UpdateItemQuantity(r.Context(), userID, itemID, req.Quantity)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "cart item not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrInvalidQuantity) {
			http.Error(w, "quantity must be greater than zero", http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := dto.CartItemResponse{
		ID:        item.ID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

func (h *cartHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	itemID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteItem(r.Context(), userID, itemID); err != nil {
		log.Printf("delete cart item error: %v", err)
		if errors.Is(err, service.ErrCartItemNotFound) {
			http.Error(w, "item id not found", http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *cartHandler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if err := h.service.DeleteCart(r.Context(), userID); err != nil {
		if errors.Is(err, service.ErrCartNotFound) {
			http.Error(w, "cart not found", http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
