package service

import (
	"context"
	"time"

	public_model "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/public/model"
	"github.com/golang-jwt/jwt"
)

type ITokenService interface {
	CreateToken(ctx context.Context, userID string, expiration int64) (string, error)
	CreateTokenPair(ctx context.Context, userID string) (*public_model.TokenModel, error)
	RefreshToken(ctx context.Context, refreshToken string) (*public_model.TokenModel, error)
}

type TokenService struct {
	JWTSecret []byte
}

func NewTokenService(jwtSecret []byte) *TokenService {
	return &TokenService{
		JWTSecret: jwtSecret,
	}
}

// CreateToken implements ITokenService.
func (t *TokenService) CreateToken(ctx context.Context, userID string, expiration int64) (string, error) {
	claims := public_model.CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration,
		},
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := unsignedToken.SignedString(t.JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// CreateTokenPair implements ITokenService.
func (t *TokenService) CreateTokenPair(ctx context.Context, userID string) (*public_model.TokenModel, error) {
	accessToken, err := t.CreateToken(ctx, userID, time.Now().Add(time.Minute*15).Unix())
	if err != nil {
		return nil, err
	}

	refreshToken, err := t.CreateToken(ctx, userID, time.Now().Add(time.Hour*24*7).Unix())
	if err != nil {
		return nil, err
	}

	tokenModel := &public_model.TokenModel{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokenModel, nil
}

// RefreshToken implements ITokenService.
func (t *TokenService) RefreshToken(ctx context.Context, refreshToken string) (*public_model.TokenModel, error) {
	claims := &public_model.CustomClaims{}

	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return t.JWTSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	tokenModel, err := t.CreateTokenPair(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	return tokenModel, nil
}

// Ensure that the service implements the interface
var _ ITokenService = (*TokenService)(nil)
