package suppfunc

import (
	"auth/internal/config"
	"auth/internal/transport/models"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateOpaqueToken(length int) (string, error) {
	b := make([]byte, length)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func GenerateJwtToken(guid string, refresh_token string, cfg *config.Config) (string, *models.CustomClaims, error) {

	claims := models.CustomClaims{
		GUID:         guid,
		RefreshToken: refresh_token,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString(cfg.SecretKey)
	if err != nil {
		return "", nil, err
	}
	return tokenString, &claims, nil
}

func ValidateJWT(tokenString string, refreshToken string, cfg *config.Config) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return cfg.SecretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("Error validation token: %v", err)
	}
	if !token.Valid {
		return "", fmt.Errorf("Token is invalid")
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && claims.RefreshToken == refreshToken {
		return claims.GUID, nil
	}

	return "", fmt.Errorf("invalid token or claims")
}
