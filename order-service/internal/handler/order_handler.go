package handler

import (
	"encoding/json"
	"net/http"

	http2 "github.com/esmrzv/microservice-shop/order-service/internal/http"
	"github.com/esmrzv/microservice-shop/order-service/internal/middleware"
	"github.com/esmrzv/microservice-shop/order-service/internal/models"
	"github.com/esmrzv/microservice-shop/order-service/internal/service"
	"github.com/google/uuid"
)

type OrderHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
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
