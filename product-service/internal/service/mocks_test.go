package service

import (
	"context"

	"github.com/esmrzv/product-service/internal/dto"
	"github.com/esmrzv/product-service/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type mockProductRepository struct {
	createProductFunc  func(context.Context, *models.Product) error
	getProductByIDFunc func(context.Context, uuid.UUID) (*models.Product, error)
	getAllProductsFunc func(context.Context) ([]models.Product, error)
	updateProductFunc  func(context.Context, *models.Product) (bool, error)
	deleteProductFunc  func(context.Context, uuid.UUID) (bool, error)
}

type mockCategoryRepository struct {
	getCategoryByIDFunc func(context.Context, uuid.UUID) (*models.Category, error)
	createCategoryFunc  func(context.Context, string) (*models.Category, error)
}

func (m *mockProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	return m.createProductFunc(ctx, product)
}

func (m *mockProductRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	return m.getProductByIDFunc(ctx, id)
}

func (m *mockProductRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	return m.getAllProductsFunc(ctx)
}

func (m *mockProductRepository) UpdateProduct(ctx context.Context, product *models.Product) (bool, error) {
	return m.updateProductFunc(ctx, product)
}

func (m *mockProductRepository) DeleteProduct(ctx context.Context, id uuid.UUID) (bool, error) {
	return m.deleteProductFunc(ctx, id)
}

func (m *mockCategoryRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	return m.getCategoryByIDFunc(ctx, id)
}

func (m *mockCategoryRepository) CreateCategory(ctx context.Context, category string) (*models.Category, error) {
	return m.createCategoryFunc(ctx, category)
}

type mockProductCache struct{}

func (m *mockProductCache) Get(
	ctx context.Context,
	id uuid.UUID,
) (*dto.ProductResponse, error) {
	return nil, redis.Nil
}

func (m *mockProductCache) Set(
	ctx context.Context,
	product *dto.ProductResponse,
) error {
	return nil
}

func (m *mockProductCache) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return nil
}
