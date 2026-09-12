package repository

import (
	"context"

	"github.com/esmrzv/cart-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepository interface {
	GetCartByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	CreateCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	GetItemsByCartID(ctx context.Context, cartID uuid.UUID) ([]*models.CartItem, error)
	AddItem(ctx context.Context, cartID uuid.UUID, productID uuid.UUID, quantity int) (*models.CartItem, error)
	UpdateItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) (*models.CartItem, error)
	DeleteItem(ctx context.Context, itemID uuid.UUID) (bool, error)
}
type cartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) CartRepository {
	return &cartRepository{
		db: db,
	}
}

func (r *cartRepository) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	query := `SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = $1`
	var cart models.Cart
	if err := r.db.QueryRow(ctx, query, userID).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt); err != nil {
		return nil, err
	}

	return &cart, nil
}
func (r *cartRepository) CreateCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	query := `INSERT INTO carts (user_id) 
			  VALUES ($1)
			  RETURNING id, user_id, created_at, updated_at`
	var cart models.Cart
	if err := r.db.QueryRow(ctx, query, userID).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt); err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetItemsByCartID(ctx context.Context, cartID uuid.UUID) ([]*models.CartItem, error) {
	query := `SELECT id, cart_id, product_id, quantity 
			  FROM cart_items WHERE cart_id = $1
			  ORDER BY id`
	rows, err := r.db.Query(ctx, query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*models.CartItem
	for rows.Next() {
		var item models.CartItem
		if err := rows.Scan(&item.ID, &item.CartID, &item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *cartRepository) AddItem(
	ctx context.Context,
	cartID uuid.UUID,
	productID uuid.UUID,
	quantity int,
) (*models.CartItem, error) {
	query := `INSERT INTO cart_items (cart_id, product_id, quantity)
			  VALUES ($1, $2, $3)
			  ON CONFLICT (cart_id, product_id)
			  DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
			  RETURNING id, cart_id, product_id, quantity`
	var item models.CartItem
	if err := r.db.QueryRow(ctx, query, cartID, productID, quantity).
		Scan(&item.ID, &item.CartID, &item.ProductID, &item.Quantity); err != nil {
		return nil, err
	}
	return &item, nil

}

func (r *cartRepository) UpdateItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) (*models.CartItem, error) {
	query := `UPDATE cart_items
			SET quantity = $1
			WHERE id = $2
			RETURNING id, cart_id, product_id, quantity
			`
	var item models.CartItem
	if err := r.db.QueryRow(ctx, query, quantity, itemID).Scan(&item.ID, &item.CartID, &item.ProductID, &item.Quantity); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) DeleteItem(ctx context.Context, itemID uuid.UUID) (bool, error) {
	query := `DELETE FROM cart_items WHERE id = $1`
	res, err := r.db.Exec(ctx, query, itemID)
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}
