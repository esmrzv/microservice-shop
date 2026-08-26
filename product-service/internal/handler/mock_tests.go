package handler

import (
	"context"

	"github.com/esmrzv/product-service/internal/dto"
	"github.com/google/uuid"
)

type mockProductService struct {
	createProductFunc  func(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error)
	getProductByIDFunc func(ctx context.Context, productID uuid.UUID) (*dto.ProductResponse, error)
	getAllProductsFunc func(ctx context.Context) ([]dto.ProductResponse, error)
	updateProductFunc  func(ctx context.Context, productID uuid.UUID, product *dto.UpdateProductRequest) error
	deleteProductFunc  func(ctx context.Context, productID uuid.UUID) error
}

func (m *mockProductService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	return m.createProductFunc(ctx, req)
}

func (m *mockProductService) GetProductByID(ctx context.Context, productID uuid.UUID) (*dto.ProductResponse, error) {
	return m.getProductByIDFunc(ctx, productID)
}

func (m *mockProductService) GetAllProducts(ctx context.Context) ([]dto.ProductResponse, error) {
	return m.getAllProductsFunc(ctx)
}

func (m *mockProductService) UpdateProduct(ctx context.Context, productID uuid.UUID, product *dto.UpdateProductRequest) error {
	return m.updateProductFunc(ctx, productID, product)
}

func (m *mockProductService) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	return m.deleteProductFunc(ctx, productID)
}
