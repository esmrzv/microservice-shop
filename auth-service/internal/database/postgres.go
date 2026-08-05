package database

import (
	"context"
	"fmt"
	"time"

	"github.com/esmrzv/microservice-shop/auth-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)
const(
	DefaultMaxConns = 10
	DefaultMinConns = 2
	MaxConnIdleTime = time.Hour
)

func NewPostges(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
						cfg.DBUser,
						cfg.DBPassword,
						cfg.DBHost,
						cfg.DBPort,
						cfg.DBName,
					)
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil{
		return nil, err
	}

	poolConfig.MaxConns = DefaultMaxConns 
	poolConfig.MinConns = DefaultMinConns
	poolConfig.MaxConnIdleTime = MaxConnIdleTime


	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil{
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil{
		pool.Close()
		return nil, err
	}


	return pool, nil

	
} 