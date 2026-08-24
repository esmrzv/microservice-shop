package models

import "time"
import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID
	CategoryID  uuid.UUID
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Category struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
