package http

import (
	"errors"
	"net/http"

	"github.com/esmrzv/microservice-shop/order-service/internal/service"
)

func GetErrors(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidQuantity):
		http.Error(w, "Invalid Quantity", http.StatusBadRequest)
	case errors.Is(err, service.ErrProductNotFound):
		http.Error(w, "Product not found", http.StatusNotFound)
	default:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
