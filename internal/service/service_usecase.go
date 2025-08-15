package service

import (
	"auth/internal/config"
	"auth/internal/repository"
	suppfunc "auth/internal/service/supp_func"
	"auth/internal/transport/models"
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type AuthServiceImpl struct {
	Repo   repository.RepositoryUC
	Logger *zap.Logger
	Config *config.Config
}

func NewAuthService(logger *zap.Logger, cfg *config.Config) (AuthServiceUC, error) {
	repo, err := repository.NewRepository(cfg, logger, context.Background())
	if err != nil {
		logger.Error("Error create repository", zap.Error(err))
		return nil, err
	}
	return &AuthServiceImpl{
		Repo:   repo,
		Logger: logger,
		Config: cfg,
	}, nil
}

func (a *AuthServiceImpl) TakeBothTokens(guid string, ctx context.Context) (*models.Tokens, error) {
	err := a.Repo.CheckTokens_AndUserExists(guid, ctx)
	if err.Error() == "Tokens exists" {
		return nil, err
	} else if err.Error() == "User doesn't exists" {
		return nil, err
	} else if err.Error() == "Tokens Already exists" {

	} else if err != nil {
		a.Logger.Error("Error check tokens.", zap.Error(err))
		return nil, err
	}

	refreshToken, err := suppfunc.GenerateOpaqueToken(32)
	if err != nil {
		a.Logger.Error("Error create refresh token.", zap.Error(err))
		return nil, err
	}

	accessToken, claims, err := suppfunc.GenerateJwtToken(guid, refreshToken, a.Config)
	if err != nil {
		a.Logger.Error("Error generate JWT access token.", zap.Error(err))
		return nil, err
	}

	err = a.Repo.InsertRefreshToken(refreshToken, claims, ctx)
	if err != nil {
		a.Logger.Error("Error insert refresh token.")
		return nil, err
	}

	return &models.Tokens{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}
func (a *AuthServiceImpl) RefreshTokens(refresh_token string, access_token string, userAgent string, ctx context.Context) (*models.Tokens, error) {
	b, err := a.Repo.CheckUserAgent(refresh_token, userAgent, ctx)
	if !b {
		return nil, errors.New("Changed User-Agent")
	}

	guid, err := suppfunc.ValidateJWT(access_token, refresh_token, a.Config)
	if err != nil {
		a.Logger.Info("Error validate JWT", zap.Error(err))
		return nil, err
	}

	err = a.Repo.CheckRefreshTokenFamily_AndCheckDeauthorized(guid, refresh_token, ctx)
	if err != nil {
		if err.Error() == "User deauthorized" || err.Error() == "Token not valid" {
			return nil, fmt.Errorf("Unauthorized access attempt: %s", err.Error())
		}
		a.Logger.Error("Error check refresh token family or check auth of user", zap.Error(err))
		return nil, err
	}

	newRefreshToken, err := suppfunc.GenerateOpaqueToken(32)
	if err != nil {
		a.Logger.Error("Error generate refresh token", zap.Error(err))
		return nil, err
	}

	newAccessToken, _, err := suppfunc.GenerateJwtToken(guid, newRefreshToken, a.Config)
	if err != nil {
		a.Logger.Error("Error generate JWT access token.", zap.Error(err))
		return nil, err
	}

	err = a.Repo.SetInvalidRefreshToken(refresh_token, ctx)
	if err != nil {
		a.Logger.Error("Error set refresh token to nvalid", zap.Error(err))
		return nil, err
	}
	return &models.Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
func (a *AuthServiceImpl) TakeGUID(user *models.User, ctx context.Context) (string, error) {
	guid, err := a.Repo.GetGUID(user, ctx)
	if err != nil {
		if err.Error() == "user deauthorized" {
			return "", err
		}
		a.Logger.Error("Error get GUID.", zap.Error(err))
		return "", err
	}
	return guid, nil
}
func (a *AuthServiceImpl) Deauthorization(accessToken string, ctx context.Context) (error) {
	token, err := jwt.ParseWithClaims(accessToken, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.Config.SecretKey, nil
	})
	if err != nil {
		return fmt.Errorf("Error validation token: %v", err)
	}
	if !token.Valid {
		return fmt.Errorf("Token is invalid")
	}
	claims, ok := token.Claims.(*models.CustomClaims)
	if ok {
		err = a.Repo.DeauthorizeByRefreshToken(claims.RefreshToken, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
