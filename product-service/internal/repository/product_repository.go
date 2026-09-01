package repository

import (
	"context"

	"github.com/esmrzv/product-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProductByID(ctx context.Context, productID uuid.UUID) (*models.Product, error)
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (bool, error)
	DeleteProduct(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (bool, error)
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
		INSERT INTO products (user_id, category_id, name, description, price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		ctx,
		query,
		product.UserID,
		product.CategoryID,
		product.Name,
		product.Description,
		product.Price,
	).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	return err
}

func (r *productRepository) GetProductByID(ctx context.Context, productID uuid.UUID) (*models.Product, error) {
	query := `SELECT id, user_id, category_id, name, description, price, created_at, updated_at FROM products WHERE id = $1`
	product := &models.Product{}

	err := r.db.QueryRow(ctx, query, productID).Scan(&product.ID, &product.UserID, &product.CategoryID, &product.Name, &product.Description, &product.Price,
		&product.CreatedAt, &product.UpdatedAt)

	return product, err

}

func (r *productRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	query := `SELECT id, name, description, price, created_at, updated_at 
				FROM products
				ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err

	}
	return products, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, product *models.Product) (bool, error) {
	query := `UPDATE products
			SET name = $2, description = $3, price = $4, updated_at = NOW()
			WHERE id = $1 AND user_id = $5`

	result, err := r.db.Exec(ctx, query, product.ID, product.Name, product.Description, product.Price, product.UserID)
	if err != nil {
		return false, err
	}
	return result.RowsAffected() > 0, nil

}

func (r *productRepository) DeleteProduct(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (bool, error) {
	query := `DELETE FROM products WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, productID, userID)
	if err != nil {
		return false, err
	}
	return result.RowsAffected() > 0, nil
}
