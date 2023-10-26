package fiber_server

import (
	fiber_handler "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/fiber/handler"
	fiber_util "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/fiber/util"
	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	App *fiber.App
}

func NewAuthFiberServer(config *fiber.Config, handler *fiber_handler.FiberServerHandler) *FiberServer {
	fiberServer := &FiberServer{App: fiber.New(*config)}
	fiberServer.setupRoutes(*handler)
	return fiberServer
}

func (f *FiberServer) Run(port string) error {
	return f.App.Listen(port)
}

func (f *FiberServer) setupRoutes(handler fiber_handler.FiberServerHandler) {
	f.App.Post("/login", func(c *fiber.Ctx) error {
		fiberCtx := &fiber_util.FiberContextImpl{
			Ctx: c,
		}
		return handler.Login(fiberCtx)
	})

	f.App.Post("/register", func(c *fiber.Ctx) error {
		fiberCtx := &fiber_util.FiberContextImpl{
			Ctx: c,
		}
		return handler.Register(fiberCtx)
	})
}
