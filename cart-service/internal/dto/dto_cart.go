


type CartRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity int `json:"quantity"`
}

type CartResponse struct {
	ID uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity int `json:"quantity"`
}

type CartItemResponse struct {
	ID uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity int `json:"quantity"`
}

type CartItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity int `json:"quantity"`
}
