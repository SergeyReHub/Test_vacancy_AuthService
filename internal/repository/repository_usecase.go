package repository

import (
	"auth/internal/config"
	"auth/internal/repository/postgres_storage"
	"auth/internal/transport/models"
	"context"
	"go.uber.org/zap"
)

type Repository struct {
	PostgresStorage postgres_storage.PostgresUC
	Logger          *zap.Logger
}

func NewRepository(cfg *config.Config, logger *zap.Logger, ctx context.Context) (RepositoryUC, error) {
	pg_storage, err := postgres_storage.New(&cfg.DB, logger, ctx)
	if err != nil {
		logger.Panic("Postgres init failed", zap.Error(err))
		return nil, err
	}

	return &Repository{
		PostgresStorage: pg_storage,
		Logger:          logger,
	}, nil
}

func (repo *Repository) GetGUID(user *models.User, ctx context.Context) (string, error) {
	// Logic to check blacklist and correct IP (webhook)
	return repo.PostgresStorage.GetGUID(user, ctx)
}
