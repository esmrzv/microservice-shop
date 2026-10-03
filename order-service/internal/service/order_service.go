package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/esmrzv/microservice-shop/order-service/internal/grpc"
	"github.com/esmrzv/microservice-shop/order-service/internal/models"
	"github.com/esmrzv/microservice-shop/order-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidQuantity = errors.New("invalid quantity")
	ErrOrderNotFound   = errors.New("order not found")
	ErrProductNotFound = errors.New("product not found")
)

type OrderService interface {
	Create(
		ctx context.Context,
		userID uuid.UUID,
		items []*models.OrderItem,
	) (*models.Order, error)
	GetByID(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*models.Order, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
}

type orderService struct {
	repo       repository.OrderRepository
	grpcClient grpc.ProductClient
}

func NewOrderService(repo repository.OrderRepository, productClient grpc.ProductClient) OrderService {
	return &orderService{repo: repo, grpcClient: productClient}
}

func (s *orderService) Create(ctx context.Context, userID uuid.UUID, items []*models.OrderItem) (*models.Order, error) {
	if len(items) == 0 {
		return nil, ErrInvalidQuantity
	}
	order := &models.Order{
		UserID: userID,
		Status: "pending",
		Items:  make([]*models.OrderItem, 0, len(items)),
	}
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
		productGrpc, err := s.grpcClient.GetProduct(ctx, item.ProductID)
		if err != nil {
			log.Printf("gRPC GetProduct error: %v", err)

			switch status.Code(err) {
			case codes.NotFound:
				return nil, ErrProductNotFound
			default:
				return nil, err
			}
		}
		log.Printf(
			"Product received: id=%s name=%s price=%f",
			productGrpc.Id,
			productGrpc.Name,
			productGrpc.Price,
		)
		item.Price = productGrpc.Price
		order.TotalPrice += item.Price * float64(item.Quantity)
		order.Items = append(order.Items, item)
		log.Printf(
			"Order item: product=%s quantity=%d price=%f total=%f",
			item.ProductID,
			item.Quantity,
			item.Price,
			order.TotalPrice,
		)

	}
	log.Printf("Creating order in repository")
	if err := s.repo.Create(ctx, order); err != nil {
		log.Printf("Repository Create error: %v", err)
		return nil, fmt.Errorf("create order: %w", err)
	}
	log.Printf("Order created: id=%s", order.ID)
	return order, nil

}

func (s *orderService) GetByID(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*models.Order, error) {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		log.Printf("Repository GetByID error: %v", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("repository GetByID error in service: %w", err)
	}
	if order.UserID != userID {
		return nil, ErrProductNotFound
	}
	return order, nil
}

func (s *orderService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	orders, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user orders: %w", err)
	}
	return orders, nil

}
