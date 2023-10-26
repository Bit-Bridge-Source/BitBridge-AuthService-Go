package jwt

import "github.com/golang-jwt/jwt"

// JWTHandler defines methods to generate and parse JWT tokens.
type JWTHandler interface {
	Generate(claims jwt.Claims) (string, error)
	Parse(tokenString string, claims jwt.Claims) (*jwt.Token, error)
}
