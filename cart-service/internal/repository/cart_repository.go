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
	UpdateItemQuantity(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, quantity int) (*models.CartItem, error)
	DeleteItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) (bool, error)
	DeleteCart(ctx context.Context, cartID uuid.UUID) error
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

func (r *cartRepository) UpdateItemQuantity(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, quantity int) (*models.CartItem, error) {
	query := `UPDATE cart_items ci
			SET quantity = $1
			FROM carts c
			WHERE ci.id = $2
			AND ci.cart_id = c.id 
			AND c.user_id = $3
			RETURNING ci.id, ci.cart_id, ci.product_id, ci.quantity
			`
	var item models.CartItem
	if err := r.db.QueryRow(ctx, query, quantity, itemID, userID).Scan(&item.ID, &item.CartID, &item.ProductID, &item.Quantity); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) DeleteItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) (bool, error) {
	query := `
    DELETE FROM cart_items
    WHERE id = $1
      AND cart_id IN (
          SELECT id
          FROM carts
          WHERE user_id = $2
      )
`
	res, err := r.db.Exec(ctx, query, itemID, userID)
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

func (r *cartRepository) DeleteCart(ctx context.Context, cartID uuid.UUID) error {
	query := `DELETE FROM carts WHERE id = $1`
	_, err := r.db.Exec(ctx, query, cartID)
	if err != nil {
		return err
	}
	return nil
}
