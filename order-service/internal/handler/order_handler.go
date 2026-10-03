package handler

import (
	"encoding/json"
	"log"
	"net/http"

	http2 "github.com/esmrzv/microservice-shop/order-service/internal/http"
	"github.com/esmrzv/microservice-shop/order-service/internal/middleware"
	"github.com/esmrzv/microservice-shop/order-service/internal/models"
	"github.com/esmrzv/microservice-shop/order-service/internal/service"
	"github.com/google/uuid"
)

type OrderHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	GetByUserID(w http.ResponseWriter, r *http.Request)
}

type orderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) OrderHandler {
	return &orderHandler{
		service: service,
	}
}
func (h *orderHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req http2.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http2.GetErrors(w, err)
		return
	}
	items := make([]*models.OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, &models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}
	order, err := h.service.Create(r.Context(), userID, items)
	if err != nil {
		http2.GetErrors(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *orderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /orders/{id} handler called")
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	orderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		log.Printf("err handler %v", err)
		http2.GetErrors(w, err)
		return
	}

	order, err := h.service.GetByID(r.Context(), userID, orderID)
	if err != nil {
		log.Printf("order handler error: %v", err)
		http2.GetErrors(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *orderHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		http2.GetErrors(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
