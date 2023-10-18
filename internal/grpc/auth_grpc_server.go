package grpcserver

import (
	"context"
	"net"

	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/service"
	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/proto/pb"
	public_model "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/public/model"
	"google.golang.org/grpc"
)

type IAuthGRPCServer interface {
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	Run() error
	InitServer(port string, listener Listener) error
}

type ServerConfig struct {
	Listener   net.Listener
	GRPCServer *grpc.Server
}

type Listener interface {
	Listen(network, address string) (net.Listener, error)
}

type AuthGRPCServer struct {
	AuthService service.IAuthService
	Middleware  grpc.UnaryServerInterceptor
	Config      ServerConfig
	pb.UnimplementedAuthServiceServer
}

// Initialize server but do not serve yet
func (s *AuthGRPCServer) InitServer(port string, listener Listener) error {
	lis, err := listener.Listen("tcp", port)
	if err != nil {
		return err
	}
	s.Config.Listener = lis
	s.Config.GRPCServer = grpc.NewServer(
		grpc.UnaryInterceptor(s.Middleware),
	)
	pb.RegisterAuthServiceServer(s.Config.GRPCServer, s)
	return nil
}

// Start serving
func (s *AuthGRPCServer) Run() error {
	return s.Config.GRPCServer.Serve(s.Config.Listener)
}

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

// Ensure AuthGRPCServer implements IAuthGRPCServer.
var _ IAuthGRPCServer = (*AuthGRPCServer)(nil)
