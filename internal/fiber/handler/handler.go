package fiber_handler

import (
	"github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/auth"
	fiber_util "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/internal/fiber/util"
	public_model "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/public/model"
	"github.com/gofiber/fiber/v2"
)

type FiberServerHandler struct {
	AuthService auth.IAuthService
}

func NewFiberServerHandler(authService auth.IAuthService) *FiberServerHandler {
	return &FiberServerHandler{AuthService: authService}
}

func (f *FiberServerHandler) Login(c fiber_util.FiberContext) error {
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

func (f *FiberServerHandler) Register(c fiber_util.FiberContext) error {
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
