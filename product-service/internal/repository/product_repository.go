package repository

import (
	"context"

	"github.com/esmrzv/product-service/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *models.Product) error
	//GetByID(ctx context.Context, productID uuid.UUID) (*models.Product, error)
	//GetAllProducts(ctx context.Context) ([]models.Product, error)
	//UpdateProduct(ctx context.Context, product *models.Product) error
	//DeleteProduct(ctx context.Context, productID uuid.UUID) error
}

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `
			INSERT INTO products (name, description, price)
			VALUES ($1, $2, $3)
			RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, product.Name, product.Description, product.Price).Scan(&product.ID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	return err

}
