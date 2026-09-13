package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []*CartItem
}

type CartItem struct {
	ID        uuid.UUID
	CartID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int
}
