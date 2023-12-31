package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/auth"
	internal_jwt "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/jwt"
	internal_time "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/time"
	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/token"
	public_model "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/public/model"
	common_crypto "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/crypto"
	common_error "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/error"
	grpc_connector "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/grpc"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/proto/pb"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Mocks
type MockTokenService struct {
	mock.Mock
}

// CreateToken mock
func (m *MockTokenService) CreateToken(ctx context.Context, userID string, duration time.Duration) (string, error) {
	args := m.Called(ctx, userID, duration)
	return args.String(0), args.Error(1)
}

// CreateTokenPair mock
func (m *MockTokenService) CreateTokenPair(ctx context.Context, userID string) (*public_model.TokenModel, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*public_model.TokenModel), args.Error(1)
}

// RefreshToken mock
func (m *MockTokenService) RefreshToken(ctx context.Context, refreshToken string) (*public_model.TokenModel, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(*public_model.TokenModel), args.Error(1)
}

// Ensure that the mock implements the interface
var _ token.ITokenService = (*MockTokenService)(nil)

type MockCrypto struct {
	mock.Mock
}

func (m *MockCrypto) CompareHashAndPassword(hashedPassword string, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func (m *MockCrypto) GenerateFromPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

// Ensure that the mock implements the interface
var _ common_crypto.ICrypto = (*MockCrypto)(nil)

type MockGrpcConnector struct {
	mock.Mock
}

func (m *MockGrpcConnector) Connect(serviceURL string) (*grpc.ClientConn, error) {
	args := m.Called(serviceURL)
	return args.Get(0).(*grpc.ClientConn), args.Error(1)
}

// Ensure that the mock implements the interface
var _ grpc_connector.IGrpcConnector = (*MockGrpcConnector)(nil)

type MockUserServiceClient struct {
	mock.Mock
	pb.UserServiceClient
}

func (m *MockUserServiceClient) CreateUser(ctx context.Context, req *pb.CreateUserRequest, opts ...grpc.CallOption) (*pb.PublicUserResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.PublicUserResponse), args.Error(1)
}

func (m *MockUserServiceClient) GetPrivateUserByIdentifier(ctx context.Context, in *pb.IdentifierRequest, opts ...grpc.CallOption) (*pb.UserResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.UserResponse), args.Error(1)
}

func (m *MockUserServiceClient) GetPublicUserByIdentifier(ctx context.Context, in *pb.IdentifierRequest, opts ...grpc.CallOption) (*pb.PublicUserResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.PublicUserResponse), args.Error(1)
}

// Ensure that the mock implements the interface
var _ pb.UserServiceClient = (*MockUserServiceClient)(nil)

type MockAuthService struct {
	mock.Mock
	token.ITokenService
}

func (m *MockAuthService) CreateToken(ctx context.Context, userID string, expiration int64) (string, error) {
	args := m.Called(ctx, userID, expiration)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) CreateTokenPair(ctx context.Context, userID string) (*public_model.TokenModel, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*public_model.TokenModel), args.Error(1)
}

type MockJWTHandler struct {
	mock.Mock
}

// Generate implements jwt.JWTHandler.
func (m *MockJWTHandler) Generate(claims jwt.Claims) (string, error) {
	args := m.Called(claims)
	return args.String(0), args.Error(1)
}

// Parse implements jwt.JWTHandler.
func (m *MockJWTHandler) Parse(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	args := m.Called(tokenString, claims)
	return args.Get(0).(*jwt.Token), args.Error(1)
}

// Ensure that the mock implements the interface
var _ internal_jwt.JWTHandler = (*MockJWTHandler)(nil)

type MockTimeSource struct {
	mock.Mock
}

// Now implements time.TimeSource.
func (*MockTimeSource) Now() time.Time {
	return time.Now()
}

// Ensure that the mock implements the interface
var _ internal_time.TimeSource = (*MockTimeSource)(nil)

func TestRegister_Success(t *testing.T) {
	// Setup mocks
	mockTokenService := new(MockTokenService)
	mockCrypto := new(MockCrypto)
	mockUserServiceClient := new(MockUserServiceClient)

	authService := auth.NewAuthService(
		mockTokenService,
		mockCrypto,
		mockUserServiceClient,
	)

	// Setup expectations
	mockUserServiceClient.On("CreateUser", mock.Anything, mock.Anything).Return(&pb.PublicUserResponse{
		Id: "test",
	}, nil)
	mockTokenService.On("CreateToken", mock.Anything, mock.Anything, mock.Anything).Return("mocked_token", nil)
	mockTokenService.On("CreateTokenPair", mock.Anything, mock.Anything).Return(&public_model.TokenModel{AccessToken: "mocked_access_token", RefreshToken: "mocked_refresh_token"}, nil)

	// Call method
	registerModel := &public_model.RegisterModel{Email: "test@test.com", Username: "test", Password: "password"}
	result, err := authService.Register(context.Background(), registerModel)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify that expected methods were called
	mockUserServiceClient.AssertExpectations(t)
}

func TestRegister_CreateUser_Failure_Unknown_Error(t *testing.T) {
	// Setup mocks
	mockTokenService := new(MockTokenService)
	mockCrypto := new(MockCrypto)
	mockUserServiceClient := new(MockUserServiceClient)

	authService := auth.NewAuthService(
		mockTokenService,
		mockCrypto,
		mockUserServiceClient,
	)

	mockUserServiceClient.On("CreateUser", mock.Anything, mock.Anything).Return((*pb.PublicUserResponse)(nil), errors.New("create user error"))
	mockTokenService.On("CreateToken", mock.Anything, mock.Anything, mock.Anything).Return("mocked_token", nil)
	mockTokenService.On("CreateTokenPair", mock.Anything, mock.Anything).Return(&public_model.TokenModel{AccessToken: "mocked_access_token", RefreshToken: "mocked_refresh_token"}, nil)

	registerModel := &public_model.RegisterModel{Email: "test@test.com", Username: "test", Password: "password"}
	result, err := authService.Register(context.Background(), registerModel)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "create user error", err.Error())

	mockUserServiceClient.AssertExpectations(t)
}

func TestRegister_CreateUser_Failure_Known_Error(t *testing.T) {
	// Setup mocks
	mockTokenService := new(MockTokenService)
	mockCrypto := new(MockCrypto)
	mockUserServiceClient := new(MockUserServiceClient)

	authService := auth.NewAuthService(
		mockTokenService,
		mockCrypto,
		mockUserServiceClient,
	)

	mockUserServiceClient.On("CreateUser", mock.Anything, mock.Anything).Return((*pb.PublicUserResponse)(nil), status.Errorf(400, "create user error"))
	mockTokenService.On("CreateToken", mock.Anything, mock.Anything, mock.Anything).Return("mocked_token", nil)
	mockTokenService.On("CreateTokenPair", mock.Anything, mock.Anything).Return(&public_model.TokenModel{AccessToken: "mocked_access_token", RefreshToken: "mocked_refresh_token"}, nil)

	registerModel := &public_model.RegisterModel{Email: "test@test.com", Username: "test", Password: "password"}
	result, err := authService.Register(context.Background(), registerModel)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "create user error", err.Error())

	mockUserServiceClient.AssertExpectations(t)
}

func TestRegister_CreateToken_Failure(t *testing.T) {
	// Setup mocks
	mockTokenService := new(MockTokenService)
	mockCrypto := new(MockCrypto)
	mockUserServiceClient := new(MockUserServiceClient)

	authService := auth.NewAuthService(
		mockTokenService,
		mockCrypto,
		mockUserServiceClient,
	)

	mockTokenService.On("CreateToken", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("create token error"))

	registerModel := &public_model.RegisterModel{Email: "test@mail", Username: "test", Password: "password"}
	result, err := authService.Register(context.Background(), registerModel)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "create token error", err.Error())

	mockUserServiceClient.AssertExpectations(t)
}

func TestRegister_CreateTokenPair_Failure(t *testing.T) {
	// Setup mocks
	mockTokenService := new(MockTokenService)
	mockCrypto := new(MockCrypto)
	mockUserServiceClient := new(MockUserServiceClient)

	authService := auth.NewAuthService(
		mockTokenService,
		mockCrypto,
		mockUserServiceClient,
	)

	// Setup expectations
	mockUserServiceClient.On("CreateUser", mock.Anything, mock.Anything).Return(&pb.PublicUserResponse{
		Id: "test",
	}, nil)
	mockTokenService.On("CreateToken", mock.Anything, mock.Anything, mock.Anything).Return("mocked_token", nil)
	mockTokenService.On("CreateTokenPair", mock.Anything, mock.Anything).Return((*public_model.TokenModel)(nil), errors.New("create token pair error"))

	// Call method
	registerModel := &public_model.RegisterModel{Email: "test@test.com", Username: "test", Password: "password"}
	result, err := authService.Register(context.Background(), registerModel)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "create token pair error", err.Error())

	// Verify that expected methods were called
	mockUserServiceClient.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	// Setup mocks
	mockTokenService := new(MockTokenService)
	mockCrypto := new(MockCrypto)
	mockUserServiceClient := new(MockUserServiceClient)

	authService := auth.NewAuthService(
		mockTokenService,
		mockCrypto,
		mockUserServiceClient,
	)

	// Setup expectations
	mockTokenService.On("CreateToken", mock.Anything, mock.Anything, mock.Anything).Return("mocked_token", nil)
	mockUserServiceClient.On("GetPrivateUserByIdentifier", mock.Anything, mock.Anything).Return(&pb.UserResponse{
		Id:   "test",
		Hash: "hashed_password",
	}, nil)
	mockCrypto.On("CompareHashAndPassword", "hashed_password", "password").Return(nil)
	mockTokenService.On("CreateTokenPair", mock.Anything, mock.Anything).Return(&public_model.TokenModel{
		AccessToken:  "mocked_access_token",
		RefreshToken: "mocked_refresh_token",
	}, nil)

	// Call method
	loginModel := &public_model.LoginModel{Email: "test@test.com", Password: "password"}
	result, err := authService.Login(context.Background(), loginModel)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify that expected methods were called
	mockTokenService.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
	mockUserServiceClient.AssertExpectations(t)
}

func TestLogin_CreateToken_Failure(t *testing.T) {
	// Setup mocks
	mockTokenService := new(MockTokenService)
	mockCrypto := new(MockCrypto)
	mockUserServiceClient := new(MockUserServiceClient)

	authService := auth.NewAuthService(
		mockTokenService,
		mockCrypto,
		mockUserServiceClient,
	)

	mockTokenService.On("CreateToken", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("create token error"))

	loginModel := &public_model.LoginModel{Email: "test@mail", Password: "password"}
	result, err := authService.Login(context.Background(), loginModel)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "create token error", err.Error())

	mockUserServiceClient.AssertExpectations(t)
}

func TestLogin_GetPrivateUserByIdentifier_Failure(t *testing.T) {
	// Setup mocks
	mockTokenService := new(MockTokenService)
	mockCrypto := new(MockCrypto)
	mockUserServiceClient := new(MockUserServiceClient)

	authService := auth.NewAuthService(
		mockTokenService,
		mockCrypto,
		mockUserServiceClient,
	)

	// Setup expectations
	mockTokenService.On("CreateToken", mock.Anything, mock.Anything, mock.Anything).Return("mocked_token", nil)
	mockUserServiceClient.On("GetPrivateUserByIdentifier", mock.Anything, mock.Anything).Return((*pb.UserResponse)(nil), errors.New("get private user error"))

	// Call method
	loginModel := &public_model.LoginModel{Email: "test@mail.com", Password: "password"}
	result, err := authService.Login(context.Background(), loginModel)

	// Assertions
	// Check error type
	serverError, ok := err.(*common_error.ServiceError)

	assert.True(t, ok)
	assert.Equal(t, common_error.Unauthorized, serverError.Code)
	assert.Equal(t, "Invalid credentials", serverError.Message)
	assert.Nil(t, result)

	// Verify that expected methods were called
	mockTokenService.AssertExpectations(t)
	mockUserServiceClient.AssertExpectations(t)
}

func TestLogin_CompareHashAndPassword_Failure(t *testing.T) {
	// Setup mocks
	mockTokenService := new(MockTokenService)
	mockCrypto := new(MockCrypto)
	mockUserServiceClient := new(MockUserServiceClient)
	mockAuthService := new(MockAuthService)

	authService := auth.NewAuthService(
		mockTokenService,
		mockCrypto,
		mockUserServiceClient,
	)

	// Setup expectations
	mockTokenService.On("CreateToken", mock.Anything, mock.Anything, mock.Anything).Return("mocked_token", nil)
	mockUserServiceClient.On("GetPrivateUserByIdentifier", mock.Anything, mock.Anything).Return(&pb.UserResponse{
		Id:   "test",
		Hash: "hashed_password",
	}, nil)
	mockCrypto.On("CompareHashAndPassword", "hashed_password", "password").Return(errors.New("compare hash and password error"))

	// Call method
	loginModel := &public_model.LoginModel{Email: "test@mail.com", Password: "password"}
	result, err := authService.Login(context.Background(), loginModel)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "compare hash and password error", err.Error())

	// Verify that expected methods were called
	mockTokenService.AssertExpectations(t)
	mockUserServiceClient.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestLogin_CreateTokenPair_Failure(t *testing.T) {
	// Setup mocks
	mockTokenService := new(MockTokenService)
	mockCrypto := new(MockCrypto)
	mockUserServiceClient := new(MockUserServiceClient)
	mockAuthService := new(MockAuthService)

	authService := auth.NewAuthService(
		mockTokenService,
		mockCrypto,
		mockUserServiceClient,
	)

	// Setup expectations
	mockTokenService.On("CreateToken", mock.Anything, mock.Anything, mock.Anything).Return("mocked_token", nil)
	mockUserServiceClient.On("GetPrivateUserByIdentifier", mock.Anything, mock.Anything).Return(&pb.UserResponse{
		Id:   "test",
		Hash: "hashed_password",
	}, nil)
	mockCrypto.On("CompareHashAndPassword", "hashed_password", "password").Return(nil)
	mockTokenService.On("CreateTokenPair", mock.Anything, mock.Anything).Return((*public_model.TokenModel)(nil), errors.New("create token pair error"))

	// Call method
	loginModel := &public_model.LoginModel{Email: "test@mail.com", Password: "password"}
	result, err := authService.Login(context.Background(), loginModel)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "create token pair error", err.Error())

	// Verify that expected methods were called
	mockTokenService.AssertExpectations(t)
	mockUserServiceClient.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestRefreshToken_Success_(t *testing.T) {
	mockJWTHandler := new(MockJWTHandler)
	mockTimeSource := &MockTimeSource{}
	svc := token.NewTokenService(mockTimeSource, mockJWTHandler)

	// Mock Parse method since RefreshToken will call it
	mockJWTHandler.On("Parse", "someValidRefreshToken", mock.Anything).Return(&jwt.Token{Valid: true}, nil)

	// Mock Generate method twice, because RefreshToken will call CreateTokenPair -> CreateToken twice
	mockJWTHandler.On("Generate", mock.Anything).Return("newAccessToken", nil).Once()
	mockJWTHandler.On("Generate", mock.Anything).Return("newRefreshToken", nil).Once()

	tokenModel, err := svc.RefreshToken(context.TODO(), "someValidRefreshToken")

	assert.NoError(t, err)

	expectedTokenModel := &public_model.TokenModel{
		AccessToken:  "newAccessToken",
		RefreshToken: "newRefreshToken",
	}

	assert.Equal(t, expectedTokenModel, tokenModel)
	mockJWTHandler.AssertExpectations(t)
}
