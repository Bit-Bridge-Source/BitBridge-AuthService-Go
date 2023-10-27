package auth

import (
	"context"
	"time"

	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/token"
	public_model "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/public/model"
	common_crypto "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/crypto"
	grpc_connector "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/grpc"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// IAuthService defines methods for user authentication services.
type IAuthService interface {
	Register(ctx context.Context, registerModel *public_model.RegisterModel) (*public_model.TokenModel, error)
	Login(ctx context.Context, loginModel *public_model.LoginModel) (*public_model.TokenModel, error)
}

// AuthService is the struct containing services and configurations for authentication.
type AuthService struct {
	TokenService             token.ITokenService                              // Handles token creation and validation
	Crypto                   common_crypto.ICrypto                            // Handles cryptographic operations
	GrpcConnector            grpc_connector.IGrpcConnector                    // Helps in connecting to other gRPC services
	UserServiceClientCreator func(conn *grpc.ClientConn) pb.UserServiceClient // Factory function to create a new UserService client
}

// NewAuthService is a constructor for creating an instance of AuthService with necessary dependencies.
func NewAuthService(
	tokenService token.ITokenService,
	crypto common_crypto.ICrypto,
	grpcConnector grpc_connector.IGrpcConnector,
	userServiceClientCreator func(conn *grpc.ClientConn) pb.UserServiceClient,
) *AuthService {
	return &AuthService{
		TokenService:             tokenService,
		Crypto:                   crypto,
		GrpcConnector:            grpcConnector,
		UserServiceClientCreator: userServiceClientCreator,
	}
}

// getGRPCClient creates and returns a new UserServiceClient for interacting with the user service.
func (authService *AuthService) getGRPCClient() (pb.UserServiceClient, error) {
	connection, err := authService.GrpcConnector.Connect("localhost:3001")
	if err != nil {
		return nil, err
	}
	return authService.UserServiceClientCreator(connection), nil
}

// createUser creates a new user by communicating with the user service.
func (authService *AuthService) createUser(ctx context.Context, client pb.UserServiceClient, registerModel *public_model.RegisterModel, token string) (string, error) {
	md := metadata.Pairs("Authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)
	userCreate := &pb.CreateUserRequest{
		Email:    registerModel.Email,
		Username: registerModel.Username,
		Password: registerModel.Password,
	}
	resp, err := client.CreateUser(ctx, userCreate)
	if err != nil {
		return "", err
	}
	return resp.GetId(), nil
}

// Register registers a new user, creates and returns a new token pair for the registered user.
func (authService *AuthService) Register(ctx context.Context, registerModel *public_model.RegisterModel) (*public_model.TokenModel, error) {
	token, err := authService.TokenService.CreateToken(ctx, "-1", time.Duration(time.Now().Add(time.Minute*15).Unix()))
	if err != nil {
		return nil, err
	}

	client, err := authService.getGRPCClient()
	if err != nil {
		return nil, err
	}

	userID, err := authService.createUser(ctx, client, registerModel, token)
	if err != nil {
		return nil, err
	}

	tokenModel, err := authService.TokenService.CreateTokenPair(ctx, userID)
	if err != nil {
		return nil, err
	}

	return tokenModel, nil
}

// Login authenticates a user, and if successful, creates and returns a new token pair for the user.
func (authService *AuthService) Login(ctx context.Context, loginModel *public_model.LoginModel) (*public_model.TokenModel, error) {
	token, err := authService.TokenService.CreateToken(ctx, "-1", time.Duration(time.Now().Add(time.Minute*15).Unix()))
	if err != nil {
		return nil, err
	}

	client, err := authService.getGRPCClient()
	if err != nil {
		return nil, err
	}

	md := metadata.Pairs("Authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	identifierRequest := &pb.IdentifierRequest{
		UserIdentifier: loginModel.Email,
	}

	user, err := client.GetPrivateUserByIdentifier(ctx, identifierRequest)
	if err != nil {
		return nil, err
	}

	err = authService.Crypto.CompareHashAndPassword(user.GetHash(), loginModel.Password)
	if err != nil {
		return nil, err
	}

	tokenModel, err := authService.TokenService.CreateTokenPair(ctx, user.GetId())
	if err != nil {
		return nil, err
	}

	return tokenModel, nil
}

// Ensure AuthService implements IAuthService.
var _ IAuthService = (*AuthService)(nil)
