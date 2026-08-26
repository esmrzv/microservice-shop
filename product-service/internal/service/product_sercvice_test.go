package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/esmrzv/product-service/internal/dto"
	"github.com/esmrzv/product-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func TestProductService_CreateProduct(t *testing.T) {
	categoryID := uuid.New()

	categoryRepo := &mockCategoryRepository{
		getCategoryByIDFunc: func(ctx context.Context, u uuid.UUID) (*models.Category, error) {
			return &models.Category{
				ID:   categoryID,
				Name: "Electronics",
			}, nil
		},
	}

	createCalled := false
	productRepo := &mockProductRepository{
		createProductFunc: func(ctx context.Context, product *models.Product) error {
			createCalled = true
			product.ID = uuid.New()
			product.CreatedAt = time.Now()
			product.UpdatedAt = time.Now()
			return nil
		},
	}
	service := NewProductService(productRepo, categoryRepo)

	req := dto.CreateProductRequest{
		CategoryID:  categoryID,
		Name:        "RTX 2060",
		Description: "Graphics card",
		Price:       250,
	}

	response, err := service.CreateProduct(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !createCalled {
		t.Fatal("expected CreateProduct to be called")
	}
	if response == nil {
		t.Fatal("expected response, got nil")
	}
	if response.Name != req.Name {
		t.Fatalf("expected name %q, got %q", req.Name, response.Name)
	}

	if response.CategoryID != req.CategoryID {
		t.Fatalf("expected category ID %v, got %v", req.CategoryID, response.CategoryID)
	}

	if response.Price != req.Price {
		t.Fatalf("expected price %.2f, got %.2f", req.Price, response.Price)
	}
	if response.Description != req.Description {
		t.Fatalf("expected description %q, got %q",
			req.Description, response.Description)
	}
}

func TestProductService_CreateProduct_CategoryNotFound(t *testing.T) {
	categoryID := uuid.New()
	categoryRepo := &mockCategoryRepository{
		getCategoryByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.Category, error) {
			return nil, pgx.ErrNoRows
		},
	}
	createCalled := false

	productRepo := &mockProductRepository{
		createProductFunc: func(ctx context.Context, product *models.Product) error {
			createCalled = true
			return nil
		},
	}
	service := NewProductService(productRepo, categoryRepo)
	req := dto.CreateProductRequest{
		CategoryID:  categoryID,
		Name:        "RTX 2060",
		Description: "Graphics card",
		Price:       250,
	}
	response, err := service.CreateProduct(context.Background(), req)

	if !errors.Is(err, ErrCategoryNotFound) {
		t.Fatalf("expected ErrCategoryNotFound, got %v", err)
	}

	if createCalled {
		t.Fatal("expected CreateProduct to be called")
	}
	if response != nil {
		t.Fatal("expected nil response")
	}
}

func TestProductService_CreateProduct_RepositoryError(t *testing.T) {
	categoryID := uuid.New()

	repoErr := errors.New("database error")

	categoryRepo := &mockCategoryRepository{
		getCategoryByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.Category, error) {
			return &models.Category{
				ID:        categoryID,
				Name:      "Electronics",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	productRepo := &mockProductRepository{
		createProductFunc: func(ctx context.Context, product *models.Product) error {
			return repoErr
		},
	}
	service := NewProductService(productRepo, categoryRepo)
	req := dto.CreateProductRequest{
		CategoryID:  categoryID,
		Name:        "RTX 2060",
		Description: "Graphics card",
		Price:       250,
	}
	response, err := service.CreateProduct(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
	if response != nil {
		t.Fatal("expected nil response")
	}

}

func TestProductService_GetProductByID(t *testing.T) {
	productID := uuid.New()
	product := &models.Product{
		ID:          productID,
		Name:        "RTX 2060",
		Description: "graphics card",
		Price:       250,
	}
	productRepo := &mockProductRepository{
		getProductByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.Product, error) {
			return product, nil

		},
	}
	service := NewProductService(productRepo, nil)
	response, err := service.GetProductByID(context.Background(), productID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if response == nil {
		t.Fatal("expected response, got nil")
	}
	if response.Name != product.Name {
		t.Fatalf("expected name %q, got %q", product.Name, response.Name)
	}
	if response.ID != product.ID {
		t.Fatalf("expected id %v, got %v", product.ID, response.ID)
	}
	if response.Description != product.Description {
		t.Fatalf("expected description %q, got %q", product.Description, response.Description)
	}
}

func TestProductService_GetProductByID_NotFound(t *testing.T) {
	productID := uuid.New()
	productRepo := &mockProductRepository{
		getProductByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.Product, error) {
			return nil, pgx.ErrNoRows
		},
	}
	service := NewProductService(productRepo, nil)
	response, err := service.GetProductByID(context.Background(), productID)
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
	if response != nil {
		t.Fatal("expected nil response")
	}

}

func TestProductService_GetProductByID_RepositoryError(t *testing.T) {
	productID := uuid.New()
	repoErr := errors.New("database error")
	productRepo := &mockProductRepository{
		getProductByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.Product, error) {
			return nil, repoErr
		},
	}
	service := NewProductService(productRepo, nil)
	response, err := service.GetProductByID(context.Background(), productID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
	if response != nil {
		t.Fatal("expected nil response")
	}

}

func TestProductService_GetAllProducts(t *testing.T) {

	products := []models.Product{
		{
			ID:          uuid.New(),
			Name:        "RTX 2060",
			Description: "graphics card",
			Price:       250,
		},
		{
			ID:          uuid.New(),
			Name:        "RTX 2080",
			Description: "graphics card2",
			Price:       2500,
		},
	}

	productRepo := &mockProductRepository{
		getAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return products, nil
		},
	}
	service := NewProductService(productRepo, nil)
	response, err := service.GetAllProducts(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(response) != len(products) {
		t.Fatalf("expected %d products, got %d", len(products), len(response))
	}
	if response[0].Name != products[0].Name {
		t.Fatalf("expected name %q, got %q", products[0].Name, response[0].Name)
	}
	if response[0].Description != products[0].Description {
		t.Fatalf("expected description %q, got %q", products[0].Description, response[0].Description)
	}
	if response[0].Price != products[0].Price {
		t.Fatalf("expected price %f, got %f", products[0].Price, response[0].Price)
	}
	if response[1].Name != products[1].Name {
		t.Fatalf("expected name %q, got %q", products[1].Name, response[1].Name)
	}
	if response[1].Description != products[1].Description {
		t.Fatalf("expected description %q, got %q", products[1].Description, response[1].Description)
	}
	if response[1].Price != products[1].Price {
		t.Fatalf("expected price %f, got %f", products[1].Price, response[1].Price)
	}
	if response[1].ID != products[1].ID {
		t.Fatalf("expected id %v, got %v", products[1].ID, response[1].ID)
	}
}

func TestProductService_GetAllProducts_RepositoryError(t *testing.T) {
	errRepo := errors.New("database error")
	productRepo := &mockProductRepository{
		getAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return nil, errRepo
		},
	}

	service := NewProductService(productRepo, nil)
	response, err := service.GetAllProducts(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if len(response) != 0 {
		t.Fatalf("expected 0 products, got %d", len(response))
	}
	if !errors.Is(err, errRepo) {
		t.Fatalf("expected repository error, got %v", err)
	}

	if response != nil {
		t.Fatal("expected nil response")
	}
}

func TestProductService_GetAllProducts_Empty(t *testing.T) {
	productRepo := &mockProductRepository{
		getAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return []models.Product{}, nil
		},
	}
	service := NewProductService(productRepo, nil)
	response, err := service.GetAllProducts(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(response) != 0 {
		t.Fatalf("expected 0 products, got %d", len(response))
	}
}

func TestProductService_UpdateProduct(t *testing.T) {
	productID := uuid.New()
	product := &dto.UpdateProductRequest{
		Name:        "RTX 2060",
		Description: "graphics card",
		Price:       250,
	}

	productRepo := &mockProductRepository{
		updateProductFunc: func(ctx context.Context, product *models.Product) (bool, error) {
			product.ID = productID
			return true, nil
		},
	}
	service := NewProductService(productRepo, nil)
	err := service.UpdateProduct(context.Background(), productID, product)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

}

func TestProductService_UpdateProduct_NotFound(t *testing.T) {

	productRepo := &mockProductRepository{
		updateProductFunc: func(ctx context.Context, product *models.Product) (bool, error) {

			return false, nil
		},
	}
	req := &dto.UpdateProductRequest{
		Name:        "RTX 2060",
		Description: "graphics card",
		Price:       250,
	}

	service := NewProductService(productRepo, nil)
	err := service.UpdateProduct(context.Background(), uuid.New(), req)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestProductService_UpdateProduct_RepositoryError(t *testing.T) {
	repoErr := errors.New("database error")
	productID := uuid.New()
	req := &dto.UpdateProductRequest{
		Name:        "RTX 2060",
		Description: "graphics card",
		Price:       250,
	}
	productRepo := &mockProductRepository{
		updateProductFunc: func(ctx context.Context, product *models.Product) (bool, error) {
			return false, repoErr
		},
	}
	service := NewProductService(productRepo, nil)
	err := service.UpdateProduct(context.Background(), productID, req)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

func TestProductService_DeleteProduct(t *testing.T) {
	productID := uuid.New()
	productRepo := &mockProductRepository{
		deleteProductFunc: func(ctx context.Context, u uuid.UUID) (bool, error) {
			return true, nil

		},
	}
	service := NewProductService(productRepo, nil)
	err := service.DeleteProduct(context.Background(), productID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProductService_DeleteProduct_RepositoryError(t *testing.T) {
	repoErr := errors.New("database error")
	productID := uuid.New()
	productRepo := &mockProductRepository{
		deleteProductFunc: func(ctx context.Context, u uuid.UUID) (bool, error) {
			return false, repoErr
		},
	}
	service := NewProductService(productRepo, nil)
	err := service.DeleteProduct(context.Background(), productID)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected database error, got %v", err)
	}
}

func TestProductService_DeleteProduct_NotFound(t *testing.T) {
	productID := uuid.New()
	productRepo := &mockProductRepository{
		deleteProductFunc: func(ctx context.Context, u uuid.UUID) (bool, error) {
			return false, nil
		},
	}
	service := NewProductService(productRepo, nil)
	err := service.DeleteProduct(context.Background(), productID)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound error, got %v", err)
	}
}
