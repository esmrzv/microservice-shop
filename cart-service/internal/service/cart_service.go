package service

import (
	"context"
	"errors"

	"github.com/esmrzv/cart-service/internal/models"
	"github.com/esmrzv/cart-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CartService interface {
	GetCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	AddItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int) (*models.CartItem, error)
	UpdateItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) (*models.CartItem, error)
	DeleteItem(ctx context.Context, itemID uuid.UUID) error
}

type cartService struct {
	repo repository.CartRepository
}

func NewCartService(repo repository.CartRepository) CartService {
	return &cartService{repo: repo}
}

func (s *cartService) GetCart(
	ctx context.Context,
	userID uuid.UUID,
) (*models.Cart, error) {
	cart, err := s.repo.GetCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			cart, err = s.repo.CreateCart(ctx, userID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	items, err := s.repo.GetItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	cart.Items = items

	return cart, nil
}

func (s *cartService) AddItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int) (*models.CartItem, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.repo.AddItem(ctx, cart.ID, productID, quantity)

}

func (s *cartService) UpdateItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) (*models.CartItem, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}
	return s.repo.UpdateItemQuantity(ctx, itemID, quantity)
}

func (s *cartService) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
	deleted, err := s.repo.DeleteItem(ctx, itemID)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.New("cart item not found")
	}
	return nil
}

func (s *cartService) DeleteCart(ctx context.Context, userID uuid.UUID) error {
	cart, err := s.repo.GetCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("cart not found")
		}
		return err
	}
	return s.repo.DeleteCart(ctx, cart.ID)
}
