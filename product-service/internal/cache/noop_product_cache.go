package cache

import (
	"context"

	"github.com/esmrzv/product-service/internal/dto"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type NoopProductCache struct{}

func NewNoopProductCache() ProductCache {
	return &NoopProductCache{}
}

func (c *NoopProductCache) Get(
	ctx context.Context,
	id uuid.UUID,
) (*dto.ProductResponse, error) {
	return nil, redis.Nil
}

func (c *NoopProductCache) Set(
	ctx context.Context,
	product *dto.ProductResponse,
) error {
	return nil
}

func (c *NoopProductCache) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return nil
}
