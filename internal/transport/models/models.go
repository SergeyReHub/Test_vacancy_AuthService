package models

import "github.com/golang-jwt/jwt/v5"

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type CustomClaims struct {
	GUID         string `json:"guid"`
	RefreshToken string `json:"refresh_token"`
	jwt.RegisteredClaims
}
