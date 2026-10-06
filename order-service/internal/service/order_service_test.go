package service

import (
	"context"
	"errors"
	"testing"

	"github.com/esmrzv/microservice-shop/order-service/internal/models"
	product "github.com/esmrzv/microservice-shop/proto/product"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeProductClient struct {
	getProductFunc func(ctx context.Context, uuid uuid.UUID) (*product.GetProductResponse, error)
}

func (c *fakeProductClient) GetProduct(ctx context.Context, productID uuid.UUID) (*product.GetProductResponse, error) {
	return c.getProductFunc(ctx, productID)
}

type fakeOrderRepository struct {
	createFunc      func(context.Context, *models.Order) error
	getByIDFunc     func(context.Context, uuid.UUID) (*models.Order, error)
	getByUserIDFunc func(context.Context, uuid.UUID) ([]*models.Order, error)

	createCalled bool
}

func (r *fakeOrderRepository) Create(ctx context.Context, order *models.Order) error {
	r.createCalled = true
	return r.createFunc(ctx, order)
}

func (r *fakeOrderRepository) GetByID(ctx context.Context, orderID uuid.UUID) (*models.Order, error) {
	return r.getByIDFunc(ctx, orderID)
}

func (r *fakeOrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	return r.getByUserIDFunc(ctx, userID)
}

func TestOrderService_Create(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()
	repoErr := errors.New("database error")
	tests := []struct {
		name             string
		quantity         int
		productPrice     float64
		productErr       error
		wantTotal        float64
		wantErr          error
		wantCreateCalled bool
		repoErr          error
	}{
		{
			name:             "success order creation",
			quantity:         2,
			productPrice:     250,
			productErr:       nil,
			wantTotal:        500,
			wantCreateCalled: true,
			repoErr:          nil,
		},
		{
			name:             "invalid quantity",
			quantity:         0,
			productPrice:     250,
			wantErr:          ErrInvalidQuantity,
			wantCreateCalled: false,
			repoErr:          nil,
		},
		{
			name:             "product not found",
			quantity:         1,
			productPrice:     250,
			productErr:       status.Error(codes.NotFound, "product not found"),
			wantErr:          ErrProductNotFound,
			wantCreateCalled: false,
			repoErr:          nil,
		},
		{
			name:             "repo error",
			quantity:         1,
			productPrice:     250,
			wantErr:          repoErr,
			repoErr:          repoErr,
			wantCreateCalled: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var savedOrder *models.Order

			orderRepository := &fakeOrderRepository{
				createFunc: func(ctx context.Context, order *models.Order) error {
					savedOrder = order
					return tt.repoErr
				},
			}

			productClient := &fakeProductClient{
				getProductFunc: func(ctx context.Context, id uuid.UUID) (*product.GetProductResponse, error) {
					if tt.productErr != nil {
						return nil, tt.productErr
					}
					return &product.GetProductResponse{
						Id:    productID.String(),
						Price: tt.productPrice,
					}, nil
				},
			}

			orderService := NewOrderService(orderRepository, productClient)
			items := []*models.OrderItem{
				{
					ProductID: productID,
					Quantity:  tt.quantity,
				},
			}
			order, err := orderService.Create(context.Background(), userID, items)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected: %v, got: %v", tt.wantErr, err)
			}

			if orderRepository.createCalled != tt.wantCreateCalled {
				t.Fatalf("expected: %v, got: %v", tt.wantCreateCalled, orderRepository.createCalled)
			}

			if tt.wantErr != nil {
				if order != nil {
					t.Fatalf("expected nil order, got: %v", order)
				}
				return
			}

			if order == nil {
				t.Fatalf("expected order, got nil")
			}
			if order.TotalPrice != tt.wantTotal {
				t.Fatalf("expected %v, got %v", tt.wantTotal, order.TotalPrice)
			}
			if savedOrder == nil {
				t.Fatalf("expected saved order, got nil")
			}
			if len(savedOrder.Items) != 1 {
				t.Fatalf("expected 1 item, got %d", len(savedOrder.Items))
			}
			item := savedOrder.Items[0]
			if item.ProductID != productID {
				t.Fatalf("expected %v, got %v", productID, item.ProductID)
			}
			if item.Quantity != tt.quantity {
				t.Fatalf("expected %v, got %v", tt.quantity, item.Quantity)
			}
			if item.Price != tt.productPrice {
				t.Fatalf("expected %v, got %v", tt.productPrice, item.Price)
			}

		})
	}

}

func TestOrderService_GetByID(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()
	wantOrder := &models.Order{
		ID:     orderID,
		UserID: userID,
		Status: "pending",
	}
	tests := []struct {
		name      string
		repoOrder *models.Order
		repoErr   error
		wantErr   error
	}{
		{
			name:      "success",
			repoOrder: wantOrder,
			repoErr:   nil,
			wantErr:   nil,
		},
		{
			name: "foreign order",
			repoOrder: &models.Order{
				ID:     orderID,
				UserID: uuid.New(),
				Status: "pending",
			},
			wantErr: ErrOrderNotFound,
			repoErr: nil,
		},
		{
			name:    "sql error",
			repoErr: pgx.ErrNoRows,
			wantErr: ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderRepository := &fakeOrderRepository{
				getByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.Order, error) {
					return tt.repoOrder, tt.repoErr
				},
			}
			orderService := NewOrderService(orderRepository, nil)
			order, err := orderService.GetByID(context.Background(), userID, orderID)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected: %v, got: %v", tt.wantErr, err)
			}
			if tt.wantErr != nil {
				if order != nil {
					t.Fatalf("expected nil order, got: %v", order)
				}
				return
			}
			if order == nil {
				t.Fatalf("expected order, got nil")
			}
			if order.ID != orderID {
				t.Fatalf("expected %v, got %v", orderID, order.ID)
			}
			if order.UserID != userID {
				t.Fatalf("expected %v, got %v", userID, order.UserID)
			}
		})
	}
}
