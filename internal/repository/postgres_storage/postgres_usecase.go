package postgres_storage

import (
	"auth/internal/config"
	"auth/internal/transport/models"
	"auth/pkg/pool_connections"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresImpl struct {
	Pool   *pgxpool.Pool
	Logger *zap.Logger
}

func New(cfg *config.DB, logger *zap.Logger, ctx context.Context) (PostgresUC, error) {
	pool, err := pool_connections.CreatePool(cfg, ctx)
	if err != nil {
		return nil, errors.New("DB postgres error. Init pool error.\n" + err.Error())
	}
	return &PostgresImpl{
		Pool:   pool,
		Logger: logger,
	}, nil
}

func (p *PostgresImpl) GetGUID(user *models.User, ctx context.Context) (string, error) {
	return "", nil
}
