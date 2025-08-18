package models

import "github.com/golang-jwt/jwt/v5"

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ClaimsJWT struct {
	GUID         string `json:"guid"`
	RefreshToken string `json:"refresh_token"`
	jwt.RegisteredClaims
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GuidResponse struct {
	Guid string `json:"guid"`
}
