package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/esmrzv/product-service/internal/models"
	"github.com/esmrzv/product-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrCategoryNotFound = errors.New("Category not found")

type CategoryService interface {
	CreateCategory(ctx context.Context, name string) (*models.Category, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

func (s *categoryService) CreateCategory(ctx context.Context, name string) (*models.Category, error) {
	category, err := s.repo.CreateCategory(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return category, nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return category, nil
}
