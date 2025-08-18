package service

import (
	"auth/backend/internal/transport/models"
	"context"
)

type AuthServiceUC interface {
	TakeBothTokens(guid string, userAgent string, ctx context.Context) (*models.Tokens, error)
	RefreshTokens(refreshToken string, accessToken string, userAgent string, ctx context.Context) (*models.Tokens, error)
	TakeGUID(user *models.User, ctx context.Context) (string, error)
	Deauthorization(accessToken string, ctx context.Context) (error)
}
