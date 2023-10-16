package public_model

import "github.com/golang-jwt/jwt"

type TokenModel struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenRefreshModel struct {
	Token string `json:"token"`
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}
