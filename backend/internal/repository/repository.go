package repository

import (
	"auth/backend/internal/transport/models"
	"context"
)

type RepositoryUC interface {
	GetGUID(user *models.User, ctx context.Context) (string, error)
	CheckRefreshTokenFamily_AndCheckDeauthorized(guid string, refreshToken string, ctx context.Context) error
	CheckTokens_AndUserExists(guid string, ctx context.Context) error
	InsertRefreshToken(refreshToken string, claims *models.ClaimsJWT, userAgent string, ctx context.Context) error
	SetInvalidRefreshToken(refreshToken string, ctx context.Context) error
	CheckUserAgent(userAgent string, refreshToken string, ctx context.Context) (bool, error)
	DeauthorizeByRefreshToken(refreshToken string, ctx context.Context) error
	CreateUser(user *models.User, ctx context.Context) (error)
}
