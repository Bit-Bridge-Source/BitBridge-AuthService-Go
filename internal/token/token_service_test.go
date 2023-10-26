package token_test

import (
	"context"
	"testing"
	"time"

	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/token"
	public_model "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/public/model"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTimeSource struct{}

func (m *MockTimeSource) Now() time.Time {
	return time.Now()
}

type MockJWTHandler struct {
	mock.Mock
}

func (m *MockJWTHandler) Generate(claims jwt.Claims) (string, error) {
	args := m.Called(claims)
	return args.String(0), args.Error(1)
}

func (m *MockJWTHandler) Parse(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	args := m.Called(tokenString, claims)
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func TestCreateToken_Success(t *testing.T) {
	timeSource := &MockTimeSource{}
	jwtHandler := new(MockJWTHandler)
	svc := token.NewTokenService(timeSource, jwtHandler)

	jwtHandler.On("Generate", mock.Anything).Return("mockToken", nil)

	token, err := svc.CreateToken(context.TODO(), "test-user", time.Minute*15)

	assert.NoError(t, err)
	assert.Equal(t, "mockToken", token)
	jwtHandler.AssertExpectations(t)
}

func TestCreateToken_Error(t *testing.T) {
	timeSource := &MockTimeSource{}
	jwtHandler := new(MockJWTHandler)
	svc := token.NewTokenService(timeSource, jwtHandler)

	jwtHandler.On("Generate", mock.Anything).Return("", assert.AnError)

	token, err := svc.CreateToken(context.TODO(), "test-user", time.Minute*15)

	assert.Error(t, err)
	assert.Equal(t, "", token)
	jwtHandler.AssertExpectations(t)
}

func TestCreateTokenPair_Success(t *testing.T) {
	timeSource := &MockTimeSource{}
	jwtHandler := new(MockJWTHandler)
	svc := token.NewTokenService(timeSource, jwtHandler)

	jwtHandler.On("Generate", mock.Anything).Return("mockToken", nil).Twice()

	tokenPair, err := svc.CreateTokenPair(context.TODO(), "test-user")

	assert.NoError(t, err)

	expectedTokenPair := &public_model.TokenModel{
		AccessToken:  "mockToken",
		RefreshToken: "mockToken",
	}

	assert.Equal(t, expectedTokenPair, tokenPair)

	jwtHandler.AssertExpectations(t)
}

func TestCreateTokenPair_Error(t *testing.T) {
	timeSource := &MockTimeSource{}
	jwtHandler := new(MockJWTHandler)
	svc := token.NewTokenService(timeSource, jwtHandler)

	jwtHandler.On("Generate", mock.Anything).Return("", assert.AnError).Once()

	tokenPair, err := svc.CreateTokenPair(context.TODO(), "test-user")

	assert.Error(t, err)
	assert.Nil(t, tokenPair)
	jwtHandler.AssertExpectations(t)
}

func TestRefreshToken_Success(t *testing.T) {
	timeSource := &MockTimeSource{}
	jwtHandler := new(MockJWTHandler)
	svc := token.NewTokenService(timeSource, jwtHandler)

	token := &jwt.Token{Valid: true}

	// Mock the Parse method to return a valid token and claims
	jwtHandler.On("Parse", mock.Anything, mock.Anything).Return(token, nil).Once()
	// Mock the Generate method to return mock tokens
	jwtHandler.On("Generate", mock.Anything).Return("newMockToken", nil).Twice()

	newTokenPair, err := svc.RefreshToken(context.TODO(), "valid-refresh-token")

	assert.NoError(t, err)

	expectedNewTokenPair := &public_model.TokenModel{
		AccessToken:  "newMockToken",
		RefreshToken: "newMockToken",
	}

	assert.Equal(t, expectedNewTokenPair, newTokenPair)

	jwtHandler.AssertExpectations(t)
}

func TestRefreshToken_Error(t *testing.T) {
	timeSource := &MockTimeSource{}
	jwtHandler := new(MockJWTHandler)
	svc := token.NewTokenService(timeSource, jwtHandler)

	jwtHandler.On("Parse", mock.Anything, mock.Anything).Return((*jwt.Token)(nil), assert.AnError)

	tokenPair, err := svc.RefreshToken(context.TODO(), "bad-refresh-token")

	assert.Error(t, err)
	assert.Nil(t, tokenPair)
	jwtHandler.AssertExpectations(t)
}

func TestRefreshToken_Invalid(t *testing.T) {
	timeSource := &MockTimeSource{}
	jwtHandler := new(MockJWTHandler)
	svc := token.NewTokenService(timeSource, jwtHandler)

	token := &jwt.Token{}
	jwtHandler.On("Parse", mock.Anything, mock.Anything).Return(token, nil)

	tokenPair, err := svc.RefreshToken(context.TODO(), "invalid token")

	assert.Error(t, err)
	assert.Nil(t, tokenPair)
	jwtHandler.AssertExpectations(t)
}

func TestCreateTokenPair_CreateTokenError_Generate(t *testing.T) {
	mockJWTHandler := new(MockJWTHandler)
	mockTimeSource := &MockTimeSource{}
	svc := token.NewTokenService(mockTimeSource, mockJWTHandler)

	// Mock Generate to return success for the first call and error for the second call
	mockJWTHandler.On("Generate", mock.Anything).Return("mockToken", nil).Once()
	mockJWTHandler.On("Generate", mock.Anything).Return("", assert.AnError).Once()

	_, err := svc.CreateTokenPair(context.TODO(), "test-user")

	assert.Error(t, err)
	mockJWTHandler.AssertExpectations(t)
}

func TestRefreshToken_CreateTokenPairError(t *testing.T) {
	mockJWTHandler := new(MockJWTHandler)
	mockTimeSource := &MockTimeSource{}
	svc := token.NewTokenService(mockTimeSource, mockJWTHandler)

	// Mock Parse to return a valid token
	mockToken := &jwt.Token{Valid: true}
	mockJWTHandler.On("Parse", mock.Anything, mock.Anything).Return(mockToken, nil).Once()

	// Mock Generate to return an error for CreateTokenPair calls
	mockJWTHandler.On("Generate", mock.Anything).Return("", assert.AnError).Once()

	_, err := svc.RefreshToken(context.TODO(), "someValidRefreshToken")

	assert.Error(t, err)
	mockJWTHandler.AssertExpectations(t)
}
