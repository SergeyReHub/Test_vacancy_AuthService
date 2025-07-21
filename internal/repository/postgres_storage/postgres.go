package postgres_storage

import (
	"auth/internal/transport/models"
	"context"
)

type PostgresUC interface {
	GetGUID(user *models.User, ctx context.Context) (string, error)
}
