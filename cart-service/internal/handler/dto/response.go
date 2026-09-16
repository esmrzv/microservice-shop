package dto

import (
	"time"

	"github.com/google/uuid"
)

type CartResponse struct {
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Items     []CartItemResponse `json:"items"`
}

type CartItemResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}
