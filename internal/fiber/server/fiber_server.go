package fiber_server

import (
	fiber_handler "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/fiber/handler"
	fiber_util "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/fiber/util"
	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	app *fiber.App
}

func NewFiberServer(config *fiber.Config) *FiberServer {
	return &FiberServer{app: fiber.New(*config)}
}

func (f *FiberServer) SetupRoutes(handler fiber_handler.FiberServerHandler) {
	f.app.Post("/login", func(c *fiber.Ctx) error {
		fiberCtx := &fiber_util.FiberContextImpl{
			Ctx: c,
		}
		return handler.Login(fiberCtx)
	})

	f.app.Post("/register", func(c *fiber.Ctx) error {
		fiberCtx := &fiber_util.FiberContextImpl{
			Ctx: c,
		}
		return handler.Register(fiberCtx)
	})
}

func (f *FiberServer) Run(port string) error {
	return f.app.Listen(port)
}
