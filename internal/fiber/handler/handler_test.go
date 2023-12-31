package handler_test

import (
	"context"
	"testing"

	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/auth"
	fiber_handler "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/fiber/handler"
	fiber_util "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/fiber/util"
	public_model "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

// Login implements service.IAuthService.
func (m *MockAuthService) Login(ctx context.Context, loginModel *public_model.LoginModel) (*public_model.TokenModel, error) {
	args := m.Called(ctx, loginModel)
	return args.Get(0).(*public_model.TokenModel), args.Error(1)
}

// Register implements service.IAuthService.
func (m *MockAuthService) Register(ctx context.Context, registerModel *public_model.RegisterModel) (*public_model.TokenModel, error) {
	args := m.Called(ctx, registerModel)
	return args.Get(0).(*public_model.TokenModel), args.Error(1)
}

// Ensure that MockAuthService implements IAuthService
var _ auth.IAuthService = &MockAuthService{}

type MockFiberContext struct {
	mock.Mock
}

// BodyParser implements fiberserver.FiberContext.
func (m *MockFiberContext) BodyParser(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

// Context implements fiberserver.FiberContext.
func (m *MockFiberContext) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

// JSON implements fiberserver.FiberContext.
func (m *MockFiberContext) JSON(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

// Ensure that MockFiberContext implements FiberContext
var _ fiber_util.FiberContext = &MockFiberContext{}

type MockFiberCtx struct {
	mock.Mock
}

func (m *MockFiberCtx) BodyParser(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockFiberCtx) JSON(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockFiberCtx) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

func TestLogin_Success(t *testing.T) {
	// Arrange
	mockFiberContext := new(MockFiberContext)
	mockAuthService := new(MockAuthService)

	mockFiberContext.On("BodyParser", mock.Anything).Return(nil)
	mockFiberContext.On("Context").Return(context.Background())
	mockFiberContext.On("JSON", mock.Anything).Return(nil)
	mockAuthService.On("Login", mock.Anything, mock.Anything).Return(&public_model.TokenModel{}, nil)

	handler := fiber_handler.NewFiberServerHandler(mockAuthService)

	// Act
	err := handler.Login(mockFiberContext)

	// Assert
	assert.Nil(t, err)

	mockFiberContext.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestLogin_Error_BadRequest(t *testing.T) {
	// Arrange
	mockFiberContext := new(MockFiberContext)
	mockAuthService := new(MockAuthService)

	mockFiberContext.On("BodyParser", mock.Anything).Return(nil)
	mockFiberContext.On("Context").Return(context.Background())
	mockAuthService.On("Login", mock.Anything, mock.Anything).Return((*public_model.TokenModel)(nil), assert.AnError)

	handler := fiber_handler.NewFiberServerHandler(mockAuthService)

	// Act
	err := handler.Login(mockFiberContext)

	// Assert
	assert.NotNil(t, err)

	mockFiberContext.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestLogin_Error_BodyParser(t *testing.T) {
	// Arrange
	mockFiberContext := new(MockFiberContext)
	mockAuthService := new(MockAuthService)

	mockFiberContext.On("BodyParser", mock.Anything).Return(assert.AnError)

	handler := fiber_handler.NewFiberServerHandler(mockAuthService)

	// Act
	err := handler.Login(mockFiberContext)

	// Assert
	assert.NotNil(t, err)

	mockFiberContext.AssertExpectations(t)
}

func TestRegister_Success(t *testing.T) {
	// Arrange
	mockFiberContext := new(MockFiberContext)
	mockAuthService := new(MockAuthService)

	mockFiberContext.On("BodyParser", mock.Anything).Return(nil)
	mockFiberContext.On("Context").Return(context.Background())
	mockFiberContext.On("JSON", mock.Anything).Return(nil)
	mockAuthService.On("Register", mock.Anything, mock.Anything).Return(&public_model.TokenModel{}, nil)

	handler := fiber_handler.NewFiberServerHandler(mockAuthService)

	// Act
	err := handler.Register(mockFiberContext)

	// Assert
	assert.Nil(t, err)

	mockFiberContext.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestRegister_Error_BadRequest(t *testing.T) {
	// Arrange
	mockFiberContext := new(MockFiberContext)
	mockAuthService := new(MockAuthService)

	mockFiberContext.On("BodyParser", mock.Anything).Return(nil)
	mockFiberContext.On("Context").Return(context.Background())
	mockAuthService.On("Register", mock.Anything, mock.Anything).Return((*public_model.TokenModel)(nil), assert.AnError)

	handler := fiber_handler.NewFiberServerHandler(mockAuthService)

	// Act
	err := handler.Register(mockFiberContext)

	// Assert
	assert.NotNil(t, err)

	mockFiberContext.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestRegister_Error_BodyParser(t *testing.T) {
	// Arrange
	mockFiberContext := new(MockFiberContext)
	mockAuthService := new(MockAuthService)

	mockFiberContext.On("BodyParser", mock.Anything).Return(assert.AnError)

	handler := fiber_handler.NewFiberServerHandler(mockAuthService)

	// Act
	err := handler.Register(mockFiberContext)

	// Assert
	assert.NotNil(t, err)

	mockFiberContext.AssertExpectations(t)
}
