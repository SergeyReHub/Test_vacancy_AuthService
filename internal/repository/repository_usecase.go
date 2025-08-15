package repository

import (
	"auth/internal/config"
	"auth/internal/repository/postgres_storage"
	"auth/internal/transport/models"
	"context"
	"errors"

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
	guid, err := repo.PostgresStorage.GetGUID(user, ctx)
	if err != nil {
		if err.Error() == "user deauthorized" {
			return "", err
		}
		repo.Logger.Error("Postgres error get GUID.", zap.Error(err))
		return "", err
	}

	return guid, nil
}

func (repo *Repository) CheckRefreshTokenFamily_AndCheckDeauthorized(guid string, refreshToken string, ctx context.Context) error {
	err := repo.PostgresStorage.GetRefreshTokenFamily_Plus_UserAuth(guid, refreshToken, ctx)
	if err != nil {
		if err.Error() == "User deauthorized" || err.Error() == "Token is not valid. Now all token family is invalid" {
			return err
		}
		repo.Logger.Error("error check refresh token family or user authorize status", zap.Error(err))
		return err
	}
	return nil
}

func (repo *Repository) CheckTokens_AndUserExists(guid string, ctx context.Context) error {
	err := repo.PostgresStorage.CheckTokens_AndUserExists(guid, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) InsertRefreshToken(refreshToken string, claims *models.CustomClaims, ctx context.Context) error {
	issued_at := claims.IssuedAt.Time
	expires_at := claims.ExpiresAt.Time
	user_guid := claims.GUID

	err := repo.PostgresStorage.InsertRefreshToken(refreshToken, expires_at, issued_at, user_guid, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) SetInvalidRefreshToken(refreshToken string, ctx context.Context) error {

	err := repo.PostgresStorage.SetInvalidRefreshToken(refreshToken, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) CheckUserAgent(userAgent string, refreshToken string, ctx context.Context) (bool, error) {

	b, err := repo.PostgresStorage.CheckUserAgent(userAgent, refreshToken, ctx)
	if err != nil {
		return false, err
	}
	return b, nil
}
func (repo *Repository) DeauthorizeByRefreshToken (refreshToken string, ctx context.Context) (error) {

	b, err := repo.PostgresStorage.DeauthorizeByRefreshToken(refreshToken, ctx)
	if err != nil {
		return err
	}
	if !b {
		return errors.New("Refresh token is not exists")
	}
	return nil
}