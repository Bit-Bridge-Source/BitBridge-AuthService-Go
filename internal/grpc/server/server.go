package grpcserver

import (
	"context"
	"net"

	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/auth"
	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/proto/pb"
	public_model "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/public/model"
	common_grpc "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/grpc"
	"google.golang.org/grpc"
)

// IAuthGRPCServer is an interface defining the authentication related methods that the GRPC server should implement.
type IAuthGRPCServer interface {
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	Run() error
	InitServer(port string, listener common_grpc.Listener) error
}

// ServerConfig holds configuration for the server such as Listener and GRPCServer.
type ServerConfig struct {
	Listener   net.Listener
	GRPCServer *grpc.Server
}

// AuthGRPCServer is a struct that embeds the services and configurations needed for the authentication server.
type AuthGRPCServer struct {
	AuthService                       auth.IAuthService             // Authentication service
	Interceptors                      []grpc.UnaryServerInterceptor // Interceptors for the GRPC server
	Config                            ServerConfig                  // Server configuration
	pb.UnimplementedAuthServiceServer                               // Embedding the unimplemented server for forward compatibility
}

// NewAuthGRPCServer is a constructor for creating an instance of AuthGRPCServer with necessary dependencies.
func NewAuthGRPCServer(authService auth.IAuthService, interceptors []grpc.UnaryServerInterceptor) *AuthGRPCServer {
	return &AuthGRPCServer{
		AuthService:  authService,
		Interceptors: interceptors,
	}
}

// InitServer initializes the server with the given port and listener but doesn’t start it.
func (s *AuthGRPCServer) InitServer(port string, listener common_grpc.Listener) error {
	lis, err := listener.Listen("tcp", port)
	if err != nil {
		return err
	}
	s.Config.Listener = lis
	s.Config.GRPCServer = grpc.NewServer(
		grpc.UnaryInterceptor(common_grpc.ChainUnaryInterceptors(s.Interceptors...)),
	)
	pb.RegisterAuthServiceServer(s.Config.GRPCServer, s)
	return nil
}

// Run starts the GRPC server, and it will run until it's stopped.
func (s *AuthGRPCServer) Run() error {
	return s.Config.GRPCServer.Serve(s.Config.Listener)
}

// Login handles the login requests, authenticating users and returning tokens.
func (s *AuthGRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	loginModel := &public_model.LoginModel{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	token, err := s.AuthService.Login(ctx, loginModel)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

// Register handles user registration requests, creating user accounts and returning tokens.
func (s *AuthGRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	registerModel := &public_model.RegisterModel{
		Email:    req.GetEmail(),
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	token, err := s.AuthService.Register(ctx, registerModel)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

// Ensuring at compile time that AuthGRPCServer implements IAuthGRPCServer interface.
var _ IAuthGRPCServer = (*AuthGRPCServer)(nil)
