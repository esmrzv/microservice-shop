package repository

import (
	"context"

	"github.com/esmrzv/product-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, name string) (*models.Category, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
}

type categoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{db: db}
}

func (cr *categoryRepository) CreateCategory(ctx context.Context, name string) (*models.Category, error) {
	query := `INSERT INTO categories(name)
				VALUES ($1)
				RETURNING id, name, created_at, updated_at
				`
	category := &models.Category{}

	err := cr.db.QueryRow(ctx, query, name).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (cr *categoryRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	query := `SELECT id, name, created_at, updated_at FROM categories WHERE id = $1`
	category := &models.Category{}
	err := cr.db.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return category, nil
}
