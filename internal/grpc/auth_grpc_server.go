package grpcserver

import (
	"context"
	"log"
	"net"

	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/service"
	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/proto/pb"
	public_model "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/public/model"
	"google.golang.org/grpc"
)

type IAuthGRPCServer interface {
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	Run(port string) error
}

type AuthGRPCServer struct {
	AuthService service.IAuthService
	Middleware  grpc.UnaryServerInterceptor
	pb.UnimplementedAuthServiceServer
}

func (s *AuthGRPCServer) Run(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %s", err.Error())
		return err
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(s.Middleware),
	)
	pb.RegisterAuthServiceServer(grpcServer, s)
	return grpcServer.Serve(lis)
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

	user, err := s.AuthService.Register(ctx, registerModel)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
	}, nil
}

// Ensure AuthGRPCServer implements IAuthGRPCServer.
var _ IAuthGRPCServer = (*AuthGRPCServer)(nil)
