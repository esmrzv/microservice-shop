package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	CategoryID  uuid.UUID `json:"categoryId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
}

type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	CategoryID  uuid.UUID `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
