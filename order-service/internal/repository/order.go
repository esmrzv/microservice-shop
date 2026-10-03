package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/esmrzv/microservice-shop/order-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, orderID uuid.UUID) (*models.Order, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
}

type orderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) OrderRepository {
	return &orderRepository{
		pool: pool,
	}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	queryOrder := `INSERT INTO orders(
                   user_id, status, total_price)
                   VALUES ($1, $2, $3)
                   RETURNING id, created_at, updated_at`
	if err := tx.QueryRow(ctx, queryOrder, order.UserID, order.Status, order.TotalPrice).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt); err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	for _, i := range order.Items {
		err = tx.QueryRow(ctx, `
			INSERT INTO order_items 
			(order_id, product_id, quantity, price)
			VALUES ($1, $2,$3,$4)
			RETURNING id
		`, order.ID, i.ProductID, i.Quantity, i.Price).Scan(&i.ID)
		if err != nil {
			return fmt.Errorf("create order_items: %w", err)
		}
		i.OrderID = order.ID
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil

}

func (r *orderRepository) GetByID(ctx context.Context, orderID uuid.UUID) (*models.Order, error) {
	order := &models.Order{}
	log.Printf("Мы на репо")

	rows, err := r.pool.Query(ctx,
		`
	SELECT o.id, o.user_id, o.status, o.total_price, o.created_at, o.updated_at,
	       oi.id, oi.product_id, oi.quantity, oi.price
	FROM orders o
	LEFT JOIN order_items oi ON oi.order_id = o.id
	WHERE o.id = $1
	`, orderID)
	if err != nil {
		log.Println("Мы на репо2")
		return nil, fmt.Errorf("query rows: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalPrice,
			&order.CreatedAt,
			&order.UpdatedAt,
			&item.ID,
			&item.ProductID,
			&item.Quantity,
			&item.Price); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		item.OrderID = order.ID
		order.Items = append(order.Items, &item)
	}
	if err := rows.Err(); err != nil {
		log.Println("МЫ все еще на репо")
		return nil, fmt.Errorf("iterate order rows: %w", err)
	}
	if order.ID == uuid.Nil {
		return nil, pgx.ErrNoRows
	}

	return order, nil
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {

	rows, err := r.pool.Query(ctx, `
	SELECT o.id, o.user_id, o.status, o.total_price, o.created_at, o.updated_at,
	       oi.id, oi.product_id, oi.quantity, oi.price
	FROM orders o
	LEFT JOIN order_items oi ON oi.order_id = o.id
	WHERE o.user_id = $1
	ORDER BY o.created_at DESC, oi.id
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("query rows: %w", err)
	}
	defer rows.Close()
	orders := make([]*models.Order, 0)
	ordersByID := make(map[uuid.UUID]*models.Order)
	for rows.Next() {
		var (
			orderID    uuid.UUID
			userID     uuid.UUID
			status     string
			totalPrice float64
			createdAt  time.Time
			updatedAt  time.Time
			item       models.OrderItem
		)

		err := rows.Scan(
			&orderID,
			&userID,
			&status,
			&totalPrice,
			&createdAt,
			&updatedAt,
			&item.ID,
			&item.ProductID,
			&item.Quantity,
			&item.Price)
		if err != nil {
			return nil, fmt.Errorf("scan user order row: %w", err)
		}
		order, exists := ordersByID[orderID]
		if !exists {
			order = &models.Order{
				ID:         orderID,
				UserID:     userID,
				Status:     status,
				TotalPrice: totalPrice,
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
			}
			ordersByID[orderID] = order
			orders = append(orders, order)
		}
		item.OrderID = order.ID
		order.Items = append(order.Items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate order rows: %w", err)
	}
	return orders, nil

}
