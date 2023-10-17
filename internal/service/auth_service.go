package service

import (
	"context"
	"time"

	public_model "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/public/model"
	common_crypto "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/crypto"
	grpc_connector "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/grpc"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type IAuthService interface {
	Register(ctx context.Context, registerModel public_model.RegisterModel) (*public_model.TokenModel, error)
	Login(ctx context.Context, loginModel public_model.LoginModel) (*public_model.TokenModel, error)
}

type AuthService struct {
	TokenService             ITokenService
	Crypto                   common_crypto.ICrypto
	GrpcConnector            grpc_connector.IGrpcConnector
	UserServiceClientCreator func(conn *grpc.ClientConn) pb.UserServiceClient
}

func NewAuthService(
	tokenService ITokenService,
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

// Separate the logic for creating a gRPC client to make it more testable.
func (authService *AuthService) getGRPCClient() (pb.UserServiceClient, error) {
	connection, err := authService.GrpcConnector.Connect("localhost:50051")
	if err != nil {
		return nil, err
	}
	return authService.UserServiceClientCreator(connection), nil
}

// Separate out the user creation logic to a new function
func (authService *AuthService) createUser(ctx context.Context, client pb.UserServiceClient, registerModel public_model.RegisterModel, token string) (string, error) {
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

func (authService *AuthService) Register(ctx context.Context, registerModel public_model.RegisterModel) (*public_model.TokenModel, error) {
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

func (authService *AuthService) Login(ctx context.Context, loginModel public_model.LoginModel) (*public_model.TokenModel, error) {
	// Create a token
	token, err := authService.TokenService.CreateToken(ctx, "-1", time.Duration(time.Now().Add(time.Minute*15).Unix()))
	if err != nil {
		return nil, err
	}

	// Create a gRPC client
	client, err := authService.getGRPCClient()
	if err != nil {
		return nil, err
	}

	// Set the token in the metadata
	md := metadata.Pairs("Authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Make a gRPC call to retrieve the user
	identifierRequest := &pb.IdentifierRequest{
		UserIdentifier: loginModel.Email,
	}

	user, err := client.GetPrivateUserByIdentifier(ctx, identifierRequest)
	if err != nil {
		return nil, err
	}

	// Compare the password
	err = authService.Crypto.CompareHashAndPassword(user.GetHash(), loginModel.Password)
	if err != nil {
		return nil, err
	}

	// Create a token pair
	tokenModel, err := authService.TokenService.CreateTokenPair(ctx, user.GetId())
	if err != nil {
		return nil, err
	}

	return tokenModel, nil
}
