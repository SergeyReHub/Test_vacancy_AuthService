package pool_connections

import (
	"auth/internal/config"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CreatePool creates a connection pool to the database.
func CreatePool(cfg *config.DB, ctx context.Context) (*pgxpool.Pool, error) {
	log.Println(cfg.ConnStr())
	config, err := pgxpool.ParseConfig(cfg.ConnStr())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database connection string: %w", err)
	}

	config.MaxConns = cfg.MaxOpenConns
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close() // Закрываем пул при ошибке
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connection pool established successfully.")

	return pool, nil
}
