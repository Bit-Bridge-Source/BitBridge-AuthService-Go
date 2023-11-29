package main

import (
	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/app"
	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/auth"
	fiber_handler "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/fiber/handler"
	fiber_server "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/fiber/server"
	grpc_server "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/grpc/server"
	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/jwt"
	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/time"
	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/token"
	common_crypto "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/crypto"
	common_fiber "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/fiber"
	common_grpc "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/grpc"
	common_vault "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/vault"
	user_pb "github.com/Bit-Bridge-Source/BitBridge-UserService-Go/proto/pb"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func main() {
	// Initialize vault client and read secret for authentication
	vaultClient, err := common_vault.NewVault("http://127.0.0.1:8200", "XZ5!Ojk88#Ox8PoM!yZhiJfHs")
	if err != nil {
		panic(err)
	}
	vaultSecret, err := vaultClient.ReadSecret("secret/data/jwt_secret")
	if err != nil {
		panic(err)
	}

	// Initialize gRPC connection to user service
	grpcConnector := &common_grpc.GrpcConnector{}
	grpcUserConnection, err := grpcConnector.Connect("localhost:3001")
	if err != nil {
		panic(err)
	}

	// Initialize gRPC client for user service
	grpUserClient := user_pb.NewUserServiceClient(grpcUserConnection)
	jwtHandler := jwt.NewSimpleJWTHandler(vaultSecret)
	tokenService := token.NewTokenService(time.NewSystemTime(), jwtHandler)
	cryptoService := common_crypto.NewCrypto()

	authService := auth.NewAuthService(tokenService, cryptoService, grpUserClient)

	fiberHandler := fiber_handler.NewFiberServerHandler(authService)
	fiberServer := fiber_server.NewAuthFiberServer(&fiber.Config{
		ErrorHandler: common_fiber.FiberErrorHandler,
	}, fiberHandler)

	var publicMethods = map[string]struct{}{"/AuthService/Login": {}, "/AuthService/Register": {}}
	authInterceptor := common_grpc.AuthUnaryInterceptor(vaultSecret, publicMethods)
	errorInterceptor := common_grpc.GRPCErrorHandler
	grpcServer := grpc_server.NewAuthGRPCServer(authService, []grpc.UnaryServerInterceptor{authInterceptor, errorInterceptor})

	app := app.NewApp(fiberServer, grpcServer)
	app.Run(":3002", ":3003")
}
