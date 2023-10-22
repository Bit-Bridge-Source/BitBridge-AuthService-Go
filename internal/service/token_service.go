package service

import (
	"context"
	"errors"
	"time"

	public_model "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/public/model"
	"github.com/golang-jwt/jwt"
)

// ITokenService defines methods for handling token operations.
type ITokenService interface {
	CreateToken(ctx context.Context, userID string, duration time.Duration) (string, error)
	CreateTokenPair(ctx context.Context, userID string) (*public_model.TokenModel, error)
	RefreshToken(ctx context.Context, refreshToken string) (*public_model.TokenModel, error)
}

// TimeSource is an interface representing a source to get the current time.
type TimeSource interface {
	Now() time.Time
}

// JWTHandler defines methods to generate and parse JWT tokens.
type JWTHandler interface {
	Generate(claims jwt.Claims) (string, error)
	Parse(tokenString string, claims jwt.Claims) (*jwt.Token, error)
}

// TokenService contains fields necessary for token operations.
type TokenService struct {
	JWTSecret []byte     // Secret key to sign the JWT tokens
	Time      TimeSource // Source to get the current time
	JWT       JWTHandler // Handler to manage JWT tokens
}

// NewTokenService initializes a new TokenService with necessary dependencies.
func NewTokenService(jwtSecret []byte, time TimeSource, jwt JWTHandler) *TokenService {
	return &TokenService{
		JWTSecret: jwtSecret,
		Time:      time,
		JWT:       jwt,
	}
}

// CreateToken generates a new JWT token with custom claims.
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

// CreateTokenPair generates a pair of access and refresh tokens.
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

// RefreshToken validates the refresh token and generates a new token pair if valid.
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

// Ensure that TokenService implements ITokenService.
var _ ITokenService = (*TokenService)(nil)
