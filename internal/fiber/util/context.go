package fiber_util

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type FiberContext interface {
	BodyParser(v interface{}) error
	JSON(v interface{}) error
	Context() context.Context
}

type FiberContextImpl struct {
	Ctx *fiber.Ctx
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
