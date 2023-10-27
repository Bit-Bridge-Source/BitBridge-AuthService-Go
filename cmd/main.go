package main

import (
	"context"

	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/app"
	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/auth"
	fiber_handler "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/fiber/handler"
	fiber_server "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/fiber/server"
	grpc_server "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/grpc/server"
	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/jwt"
	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/time"
	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/token"
	common_crypto "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/crypto"
	grpc_connector "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/grpc"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/proto/pb"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func main() {

	jwtHandler := jwt.NewSimpleJWTHandler([]byte("secret"))
	tokenService := token.NewTokenService(time.NewSystemTime(), jwtHandler)
	cryptoService := common_crypto.NewCrypto()
	grpcConnector := &grpc_connector.GrpcConnector{}

	authService := auth.NewAuthService(tokenService, cryptoService, grpcConnector, func(conn *grpc.ClientConn) pb.UserServiceClient {
		return pb.NewUserServiceClient(conn)
	})

	fiberHandler := fiber_handler.NewFiberServerHandler(authService)
	fiberServer := fiber_server.NewAuthFiberServer(&fiber.Config{}, fiberHandler)

	grpcServer := grpc_server.NewAuthGRPCServer(authService, func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	})

	app := app.NewApp(fiberServer, grpcServer)
	app.Run(":3002", ":3003")
}
