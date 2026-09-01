package cache

import (
	"context"
	"fmt"

	"github.com/esmrzv/product-service/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, err
	}

	return client, nil
}
