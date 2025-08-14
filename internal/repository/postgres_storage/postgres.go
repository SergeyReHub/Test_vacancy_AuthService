package postgres_storage

import (
	"auth/internal/transport/models"
	"context"
	"time"
)

type PostgresUC interface {
	GetGUID(user *models.User, ctx context.Context) (string, error)
	GetRefreshTokenFamily_Plus_UserAuth(guid string, refreshToken string, ctx context.Context) error
	CheckTokens_AndUserExists(guid string, ctx context.Context) error
	InsertRefreshToken(refreshToken string, expires_at time.Time, issued_at time.Time, user_guid string, ctx context.Context) error
	SetInvalidRefreshToken(refreshToken string, ctx context.Context) error
	CheckUserAgent(userAgent string, refreshToken string, ctx context.Context) (bool, error)
}
