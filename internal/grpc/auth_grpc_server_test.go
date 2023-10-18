package grpcserver_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	grpcserver "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/grpc"
	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/proto/pb"
	public_model "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockListener struct{}

func (m *MockListener) Listen(network, address string) (net.Listener, error) {
	return nil, nil // Mock implementation
}

type MockListenerWithError struct{}

func (m *MockListenerWithError) Listen(network, address string) (net.Listener, error) {
	return nil, fmt.Errorf("mock error")
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, model *public_model.LoginModel) (*public_model.TokenModel, error) {
	args := m.Called(ctx, model)
	return args.Get(0).(*public_model.TokenModel), args.Error(1)
}

func (m *MockAuthService) Register(ctx context.Context, model *public_model.RegisterModel) (*public_model.TokenModel, error) {
	args := m.Called(ctx, model)
	return args.Get(0).(*public_model.TokenModel), args.Error(1)
}

func TestAuthGRPCServer_InitServer_Success(t *testing.T) {
	s := &grpcserver.AuthGRPCServer{}
	err := s.InitServer(":5000", &MockListener{})

	assert.Nil(t, err)
	assert.NotNil(t, s.Config.GRPCServer)
}

func TestAuthGRPCServer_InitServer_Error(t *testing.T) {
	s := &grpcserver.AuthGRPCServer{}
	err := s.InitServer(":5000", &MockListenerWithError{})

	assert.NotNil(t, err)
}

func TestAuthGRPCServer_Run_Success(t *testing.T) {
	// Mock Listener that we can close to stop the server
	lis, err := net.Listen("tcp", ":0") // :0 picks a random open port
	assert.Nil(t, err)

	s := &grpcserver.AuthGRPCServer{
		Config: grpcserver.ServerConfig{
			Listener:   lis,
			GRPCServer: grpc.NewServer(),
		},
	}
	pb.RegisterAuthServiceServer(s.Config.GRPCServer, s)

	done := make(chan error)
	go func() {
		done <- s.Run()
	}()

	// Stop the server
	s.Config.GRPCServer.Stop()

	// Check that Run() returned the expected error after server was stopped
	err = <-done
	assert.Equal(t, "grpc: the server has been stopped", err.Error())
}

// Test Login method
func TestAuthGRPCServer_Login_Success(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockAuthService.On("Login", mock.Anything, mock.Anything).Return(&public_model.TokenModel{
		AccessToken:  "expected_access_token",
		RefreshToken: "expected_refresh_token",
	}, nil)

	s := &grpcserver.AuthGRPCServer{AuthService: mockAuthService}

	req := &pb.LoginRequest{
		Email:    "test@test.com",
		Password: "password",
	}
	resp, err := s.Login(context.TODO(), req)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, "expected_access_token", resp.GetAccessToken())
	assert.Equal(t, "expected_refresh_token", resp.GetRefreshToken())

	mockAuthService.AssertExpectations(t)
}

// Test Login method with an expected error
func TestAuthGRPCServer_Login_Error(t *testing.T) {
	mockAuthService := new(MockAuthService)
	expectedError := fmt.Errorf("login failed")
	mockAuthService.On("Login", mock.Anything, mock.Anything).Return((*public_model.TokenModel)(nil), expectedError)

	s := &grpcserver.AuthGRPCServer{AuthService: mockAuthService}

	req := &pb.LoginRequest{
		Email:    "test@test.com",
		Password: "wrong_password",
	}
	resp, err := s.Login(context.TODO(), req)

	// Assertions
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)

	mockAuthService.AssertExpectations(t)
}

// Test Register method
func TestAuthGRPCServer_Register_Success(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockAuthService.On("Register", mock.Anything, mock.Anything).Return(&public_model.TokenModel{
		AccessToken:  "expected_access_token",
		RefreshToken: "expected_refresh_token",
	}, nil)

	s := &grpcserver.AuthGRPCServer{AuthService: mockAuthService}

	req := &pb.RegisterRequest{
		Email:    "test@test.com",
		Username: "username",
		Password: "password",
	}
	resp, err := s.Register(context.TODO(), req)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, "expected_access_token", resp.GetAccessToken())
	assert.Equal(t, "expected_refresh_token", resp.GetRefreshToken())

	mockAuthService.AssertExpectations(t)
}

// Test Register method with an expected error
func TestAuthGRPCServer_Register_Error(t *testing.T) {
	mockAuthService := new(MockAuthService)
	expectedError := fmt.Errorf("register failed")
	mockAuthService.On("Register", mock.Anything, mock.Anything).Return((*public_model.TokenModel)(nil), expectedError)

	s := &grpcserver.AuthGRPCServer{AuthService: mockAuthService}

	req := &pb.RegisterRequest{
		Email:    "test@test.com",
		Username: "username",
		Password: "password",
	}
	resp, err := s.Register(context.TODO(), req)

	// Assertions
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)

	mockAuthService.AssertExpectations(t)
}
