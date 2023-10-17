package service

import (
	"context"
	"errors"
	"time"

	public_model "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/public/model"
	"github.com/golang-jwt/jwt"
)

type ITokenService interface {
	CreateToken(ctx context.Context, userID string, duration time.Duration) (string, error)
	CreateTokenPair(ctx context.Context, userID string) (*public_model.TokenModel, error)
	RefreshToken(ctx context.Context, refreshToken string) (*public_model.TokenModel, error)
}

type TimeSource interface {
	Now() time.Time
}

type JWTHandler interface {
	Generate(claims jwt.Claims) (string, error)
	Parse(tokenString string, claims jwt.Claims) (*jwt.Token, error)
}

type TokenService struct {
	JWTSecret []byte
	Time      TimeSource
	JWT       JWTHandler
}

// NewTokenService creates a new TokenService.
func NewTokenService(jwtSecret []byte, time TimeSource, jwt JWTHandler) *TokenService {
	return &TokenService{
		JWTSecret: jwtSecret,
		Time:      time,
		JWT:       jwt,
	}
}

// CreateToken implements ITokenService.
func (t *TokenService) CreateToken(ctx context.Context, userID string, duration time.Duration) (string, error) {
	expiration := t.Time.Now().Add(duration).Unix()

	claims := public_model.CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration,
		},
	}

	tokenString, err := t.JWT.Generate(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// CreateTokenPair implements ITokenService.
func (t *TokenService) CreateTokenPair(ctx context.Context, userID string) (*public_model.TokenModel, error) {
	accessToken, err := t.CreateToken(ctx, userID, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, err := t.CreateToken(ctx, userID, 24*7*time.Hour)
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

	token, err := t.JWT.Parse(refreshToken, claims)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	tokenModel, err := t.CreateTokenPair(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	return tokenModel, nil
}

// Ensure that the service implements the interface
var _ ITokenService = (*TokenService)(nil)
