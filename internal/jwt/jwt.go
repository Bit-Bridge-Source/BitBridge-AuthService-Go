package jwt

import "github.com/golang-jwt/jwt"

// JWTHandler defines methods to generate and parse JWT tokens.
type JWTHandler interface {
	Generate(claims jwt.Claims) (string, error)
	Parse(tokenString string, claims jwt.Claims) (*jwt.Token, error)
}

type SimpleJWTHandler struct {
	SigningKey []byte
}

// NewSimpleJWTHandler initializes a new SimpleJWTHandler with the given signing key.
func NewSimpleJWTHandler(signingKey []byte) *SimpleJWTHandler {
	return &SimpleJWTHandler{SigningKey: signingKey}
}

// Generate implements JWTHandler.
func (s *SimpleJWTHandler) Generate(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.SigningKey)
}

// Parse implements JWTHandler.
func (s *SimpleJWTHandler) Parse(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.SigningKey, nil
	})
}

// Ensure SimpleJWTHandler implements JWTHandler.
var _ JWTHandler = (*SimpleJWTHandler)(nil)
