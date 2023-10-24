package app

import (
	fiber_server "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/fiber/server"
	grpc_server "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/grpc/server"
	common_listener "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/listener"
)

type App struct {
	FiberServer *fiber_server.FiberServer
	GRPCServer  *grpc_server.AuthGRPCServer
}

func NewApp(fiberServer *fiber_server.FiberServer, gRPCServer *grpc_server.AuthGRPCServer) *App {
	return &App{
		FiberServer: fiberServer,
		GRPCServer:  gRPCServer,
	}
}

func (app *App) Run(httpPort string, gRPCPort string) {
	// Channels to collect errors from the servers
	httpErrChan := make(chan error)
	grpcErrChan := make(chan error)

	// Run Fiber server
	go func() {
		httpErrChan <- app.FiberServer.Run(httpPort)
	}()

	// Run gRPC server
	go func() {
		grpcErrChan <- app.GRPCServer.InitServer(gRPCPort, &common_listener.DefaultListener{})
		grpcErrChan <- app.GRPCServer.Run()
	}()

	// Wait for errors from the servers
	select {
	case err := <-httpErrChan:
		panic(err)
	case err := <-grpcErrChan:
		panic(err)
	}
}
