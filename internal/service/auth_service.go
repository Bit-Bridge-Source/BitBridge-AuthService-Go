package service

import (
	"context"
	"time"

	public_model "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/public/model"
	common_crypto "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/crypto"
	grpc_connector "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/grpc"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/proto/pb"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
)

type IAuthService interface {
	Register(ctx context.Context, registerModel public_model.RegisterModel) (*public_model.TokenModel, error)
	Login(ctx context.Context, loginModel public_model.LoginModel) (*public_model.TokenModel, error)
	CreateTokenPair(ctx context.Context, userID string) (*public_model.TokenModel, error)
	CreateToken(ctx context.Context, userID string, expiration int64) (string, error)
	RefreshToken(ctx context.Context, refreshToken string) (*public_model.TokenModel, error)
}

type AuthService struct {
	JWTSecret     []byte
	Crypto        common_crypto.ICrypto
	GrpcConnector grpc_connector.IGrpcConnector
}

func NewAuthService(jwtSecret []byte, crypto common_crypto.ICrypto, grpcConnector grpc_connector.IGrpcConnector) *AuthService {
	return &AuthService{
		JWTSecret:     jwtSecret,
		Crypto:        crypto,
		GrpcConnector: grpcConnector,
	}
}

func (authService *AuthService) Register(ctx context.Context, registerModel public_model.RegisterModel) (*public_model.TokenModel, error) {
	token, err := authService.CreateToken(ctx, "-1", time.Now().Add(time.Minute*15).Unix())
	if err != nil {
		return nil, err
	}

	connection, err := authService.GrpcConnector.Connect("localhost:50051")
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	md := metadata.Pairs("Authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	userCreate := &pb.CreateUserRequest{
		Email:    registerModel.Email,
		Username: registerModel.Username,
		Password: registerModel.Password,
	}

	resp, err := client.CreateUser(ctx, userCreate)
	if err != nil {
		return nil, err
	}

	tokenModel, err := authService.CreateTokenPair(ctx, resp.GetId())
	if err != nil {
		return nil, err
	}

	return tokenModel, nil
}

func (authService *AuthService) Login(ctx context.Context, loginModel public_model.LoginModel) (*public_model.TokenModel, error) {
	token, err := authService.CreateToken(ctx, "-1", time.Now().Add(time.Minute*15).Unix())
	if err != nil {
		return nil, err
	}

	connection, err := authService.GrpcConnector.Connect("localhost:50051")
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

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

	tokenModel, err := authService.CreateTokenPair(ctx, user.GetId())
	if err != nil {
		return nil, err
	}

	return tokenModel, nil
}

func (authService *AuthService) CreateToken(ctx context.Context, userID string, expiration int64) (string, error) {
	claims := public_model.CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration,
		},
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := unsignedToken.SignedString(authService.JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (authService *AuthService) CreateTokenPair(ctx context.Context, userID string) (*public_model.TokenModel, error) {
	accessToken, err := authService.CreateToken(ctx, userID, time.Now().Add(time.Minute*15).Unix())
	if err != nil {
		return nil, err
	}

	refreshToken, err := authService.CreateToken(ctx, userID, time.Now().Add(time.Hour*24*7).Unix())
	if err != nil {
		return nil, err
	}

	tokenModel := &public_model.TokenModel{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokenModel, nil
}

func (authService *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*public_model.TokenModel, error) {
	claims := &public_model.CustomClaims{}

	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return authService.JWTSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	tokenModel, err := authService.CreateTokenPair(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	return tokenModel, nil
}
