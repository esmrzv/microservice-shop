package service

import (
	"context"
	"fmt"

	"github.com/esmrzv/product-service/internal/dto"
	"github.com/esmrzv/product-service/internal/models"
	"github.com/esmrzv/product-service/internal/repository"
)


type ProductService interface{
	CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.CreateProductResponse,error)
}



type productService struct {
	repo repository.ProductRepository
}


func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}


func (s *productService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.CreateProductResponse, error){
	product := &models.Product{
		Name: req.Name,
		Description: req.Description,
		Price: req.Price,

	}

	err := s.repo.CreateProduct(ctx, product)
	if err != nil{
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	response := &dto.CreateProductResponse{
		ID: product.ID,
		Name: product.Name,
		Description: product.Description,
		Price: product.Price,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
	return response, nil

}