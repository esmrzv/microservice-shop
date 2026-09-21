package grpc

import (
	"context"

	"github.com/esmrzv/microservice-shop/proto/product"
	"github.com/google/uuid"
)

type ProductClient interface {
	GetProduct(ctx context.Context, productID uuid.UUID) (*product.GetProductResponse, error)
}

type productClient struct {
	client product.ProductServiceClient
}

func NewProductClient(client product.ProductServiceClient) ProductClient {
	return &productClient{
		client: client,
	}
}

func (c *productClient) GetProduct(ctx context.Context, productID uuid.UUID) (*product.GetProductResponse, error) {
	req := &product.GetProductRequest{
		ProductId: productID.String(),
	}
	return c.client.GetProduct(ctx, req)
}
