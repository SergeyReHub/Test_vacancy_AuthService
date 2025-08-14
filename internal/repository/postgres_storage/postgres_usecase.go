package postgres_storage

import (
	"auth/internal/config"
	"auth/internal/transport/models"
	"auth/pkg/pool_connections"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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
	pool := p.Pool
	conn, err := pool.Acquire(ctx)
	if err != nil {
		p.Logger.Error("Error aquire pool", zap.Error(err))
		return "", err
	}
	defer pool.Close()

	row := conn.QueryRow(ctx, "SELECT FROM users WHERE username=$1 AND password=$2", user.Username, user.Password)

	var guid string
	var deauthorized bool
	err = row.Scan(&guid, &deauthorized)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", errors.New("No users with these claims")
		}
		p.Logger.Error("Error scan values", zap.Error(err))
		return "", err
	}

	if deauthorized {
		return "", errors.New("user deauthorized")
	}
	return guid, nil
}
func (p *PostgresImpl) GetRefreshTokenFamily_Plus_UserAuth(guid string, refreshToken string, ctx context.Context) error {
	pool := p.Pool
	conn, err := pool.Acquire(ctx)
	if err != nil {
		p.Logger.Error("Error aquire pool", zap.Error(err))
		return err
	}
	defer pool.Close()

	row := conn.QueryRow(ctx, `
		SELECT u.deauthorized, r.valid
		FROM refresh_tokens r
		JOIN users u
		WHERE r.token = $1 AND u.guid = $2
		LIMIT 1`, refreshToken, guid)

	var deauthorized, valid bool
	err = row.Scan(&deauthorized, &valid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("Token is not valid. Now all token family is invalid")
		}
		return err
	}
	if deauthorized {
		return errors.New("User deauthorized")
	}
	if !valid {
		//Add query that invalid all token family
		row = conn.QueryRow(ctx, `
		UPDATE refresh_tokens SET valide = false 
		WHERE user_guid=$1`, guid)
		return errors.New("Token is not valid. Now all token family is invalid")
	}
	return nil
}
func (p *PostgresImpl) CheckTokens_AndUserExists(guid string, ctx context.Context) error {
	pool := p.Pool
	conn, err := pool.Acquire(ctx)
	if err != nil {
		p.Logger.Error("Error aquire pool", zap.Error(err))
		return err
	}
	defer pool.Close()

	row := conn.QueryRow(ctx, `SELECT 
    EXISTS(SELECT 1 FROM refresh_tokens WHERE user_guid = $1) AS token_exists,
    EXISTS(SELECT 1 FROM users WHERE guid = $1) AS user_exists`, guid)

	var token_exists, user_exists bool
	err = row.Scan(&token_exists)
	if err != nil {
		p.Logger.Error("Error scan values", zap.Error(err))
		return err
	}
	if !user_exists {
		return errors.New("User doesn't exists")
	}
	if token_exists {
		return errors.New("Tokens Already exists")
	}
	return nil
}

func (p *PostgresImpl) SetInvalidRefreshToken(refreshToken string, ctx context.Context) error {
	pool := p.Pool
	conn, err := pool.Acquire(ctx)
	if err != nil {
		p.Logger.Error("Error aquire pool.", zap.Error(err))
		return err
	}
	defer pool.Close()

	hashedRefreshToken, err := hashRefreshToken(refreshToken)
	if err != nil {
		p.Logger.Error("Error bcrypt refresh token", zap.Error(err))
	}

	res, err := conn.Exec(ctx, "UPDATE refresh_tokens SET valid = false WHERE token = $1", hashedRefreshToken)
	if err != nil {
		p.Logger.Error("Error update token to set valide=false.", zap.Error(err), zap.String("res", res.String()))
	}

	return nil
}

func (p *PostgresImpl) InsertRefreshToken(refreshToken string, expires_at time.Time, issued_at time.Time, user_guid string, ctx context.Context) error {
	pool := p.Pool
	conn, err := pool.Acquire(ctx)
	if err != nil {
		p.Logger.Error("Error aquire pool.", zap.Error(err))
		return err
	}
	defer pool.Close()

	hashedRefreshToken, err := hashRefreshToken(refreshToken)
	if err != nil {
		p.Logger.Error("Error bcrypt refresh token", zap.Error(err))
	}

	res, err := conn.Exec(ctx, "INSERT INTO refresh_tokens (token, user_guid, issued_at, expires_at) VALUES ($1, $2, $3, $4, $5)", hashedRefreshToken, user_guid, expires_at, issued_at)
	if err != nil {
		p.Logger.Error("Error insert token.", zap.Error(err), zap.String("res", res.String()))
	}

	return nil
}

func (p *PostgresImpl) CheckUserAgent(userAgent string, refreshToken string, ctx context.Context) (bool, error) {
	pool := p.Pool
	conn, err := pool.Acquire(ctx)
	if err != nil {
		p.Logger.Error("Error aquire pool.", zap.Error(err))
		return false, err
	}
	defer pool.Close()

	row := conn.QueryRow(ctx, "SELECT EXISTS(SELECT FROM refresh_tokens WHERE token=$1 AND user_agent=$2)", refreshToken, userAgent)
	var b bool
	err = row.Scan(&b)
	if err != nil {
		p.Logger.Error("Error scan values", zap.Error(err))
		return false, err
	}

	return b, nil
}

func hashRefreshToken(refresh_token string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(refresh_token), 14)
	return string(bytes), err
}
