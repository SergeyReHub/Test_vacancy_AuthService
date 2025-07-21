package models

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
