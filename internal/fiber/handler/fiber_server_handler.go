package fiber_handler

import (
	"context"

	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/service"
	public_model "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/public/model"
	"github.com/gofiber/fiber/v2"
)

type FiberContext interface {
	BodyParser(v interface{}) error
	JSON(v interface{}) error
	Context() context.Context
}

type FiberContextImpl struct {
	Ctx FiberContext
}

func (f *FiberContextImpl) BodyParser(v interface{}) error {
	return f.Ctx.BodyParser(v)
}

func (f *FiberContextImpl) JSON(v interface{}) error {
	return f.Ctx.JSON(v)
}

func (f *FiberContextImpl) Context() context.Context {
	return f.Ctx.Context()
}

type FiberServerHandler struct {
	AuthService service.IAuthService
}

func NewFiberServerHandler(authService service.IAuthService) *FiberServerHandler {
	return &FiberServerHandler{AuthService: authService}
}

func (f *FiberServerHandler) Login(c FiberContext) error {
	loginModel := public_model.LoginModel{}
	if err := c.BodyParser(&loginModel); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	token, err := f.AuthService.Login(c.Context(), &loginModel)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(token)
}

func (f *FiberServerHandler) Register(c FiberContext) error {
	registerModel := public_model.RegisterModel{}
	if err := c.BodyParser(&registerModel); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	token, err := f.AuthService.Register(c.Context(), &registerModel)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(token)
}
