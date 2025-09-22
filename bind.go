package common

import "github.com/gofiber/fiber/v2"

// BindAndValidate : Bind and validate request body
func BindAndValidate[T any](ctx *fiber.Ctx) (req T, err error) {
	if err = ctx.BodyParser(&req); err != nil {
		return req, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	err = Validator().Struct(req)

	return
}
