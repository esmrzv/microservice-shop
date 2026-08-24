package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/esmrzv/product-service/internal/dto"
	"github.com/esmrzv/product-service/internal/models"
	"github.com/esmrzv/product-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrProductNotFound = errors.New("product not found")

type ProductService interface {
	CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProductByID(ctx context.Context, productID uuid.UUID) (*dto.ProductResponse, error)
	GetAllProducts(ctx context.Context) ([]dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, productID uuid.UUID, product *dto.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, productID uuid.UUID) error
}

type productService struct {
	repo         repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

func NewProductService(repo repository.ProductRepository, cr repository.CategoryRepository) ProductService {
	return &productService{
		repo:         repo,
		categoryRepo: cr,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := &models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}
	_, err := s.categoryRepo.GetCategoryByID(ctx, product.CategoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	err = s.repo.CreateProduct(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	response := &dto.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
	return response, nil

}

func (s *productService) GetProductByID(ctx context.Context, productID uuid.UUID) (*dto.ProductResponse, error) {

	product, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	response := &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
	return response, nil

}

func (s *productService) GetAllProducts(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := s.repo.GetAllProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	response := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		response[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}

	}
	return response, nil

}

func (s *productService) UpdateProduct(ctx context.Context, productID uuid.UUID, req *dto.UpdateProductRequest) error {
	product := &models.Product{
		ID:          productID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	updated, err := s.repo.UpdateProduct(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	if !updated {
		return ErrProductNotFound
	}
	return nil

}

func (s *productService) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	deleted, err := s.repo.DeleteProduct(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	if !deleted {
		return fmt.Errorf("could not delete product: %w", ErrProductNotFound)
	}
	return nil
}
