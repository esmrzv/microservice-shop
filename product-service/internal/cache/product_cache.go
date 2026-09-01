package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/esmrzv/product-service/internal/dto"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ProductCache interface {
	Get(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error)
	Set(ctx context.Context, product *dto.ProductResponse) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type productCache struct {
	client *redis.Client
}

func NewProductCache(client *redis.Client) ProductCache {
	return &productCache{
		client: client,
	}
}

func (p *productCache) Get(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	data, err := p.client.Get(ctx, id.String()).Result()
	if err != nil {
		return nil, err
	}

	var product dto.ProductResponse
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		return nil, fmt.Errorf("failed to unmarshal product: %w", err)
	}
	log.Printf("redis GET: %s", productKey(id))
	return &product, nil
}

func (p *productCache) Set(ctx context.Context, product *dto.ProductResponse) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}
	err = p.client.Set(
		ctx,
		productKey(product.ID),
		data,
		time.Hour*10,
	).Err()
	if err != nil {
		return fmt.Errorf("failed to cache product: %w", err)
	}
	return nil
}

func (p *productCache) Delete(ctx context.Context, id uuid.UUID) error {
	if err := p.client.Del(ctx, productKey(id)).Err(); err != nil {
		return fmt.Errorf("failed to delete cached product: %w", err)
	}

	return nil
}

func productKey(id uuid.UUID) string {
	return fmt.Sprintf("product:%s", id)
}
