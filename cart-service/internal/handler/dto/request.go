package dto

import "github.com/google/uuid"

type AddItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}
