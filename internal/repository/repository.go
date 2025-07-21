package repository

import (
	"auth/internal/transport/models"
	"context"
)

type RepositoryUC interface {
	GetGUID(user *models.User, ctx context.Context) (string, error)
}
