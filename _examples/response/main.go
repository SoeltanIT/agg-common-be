package main

import (
	"net/http"

	common "github.com/SoeltanIT/agg-common-be"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Send success response with data
	app.Get("/success", func(ctx *fiber.Ctx) error {
		data := fiber.Map{
			"user": "John Doe",
		}
		return common.Response().SetMessage("Success fetch user").SetData(data).Send(ctx)
	})

	// Send success response with data and custom status
	app.Get("/success-custom", func(ctx *fiber.Ctx) error {
		newData := fiber.Map{
			"user": "John Doe",
		}

		return common.Response().SetMessage("Success create user").SetData(newData, http.StatusCreated).Send(ctx)
	})

	// Send success response with pagination
	app.Get("/success-pagination", func(ctx *fiber.Ctx) error {
		p := common.NewPaginationParams(ctx)

		data, total := func() (fiber.Map, int64) {
			return fiber.Map{
				"user": "John Doe",
			}, 50
		}()

		paginationResponse := p.GetPaginationResponse(ctx.Request(), total)

		return common.Response().
			SetMessage("Success fetch user with pagination").
			SetData(data).
			SetPagination(paginationResponse).
			Send(ctx)
	})

	// Send error response
	app.Get("/error", func(ctx *fiber.Ctx) error {
		err := func() error {
			return common.ErrPlayerNotFound
		}()

		return common.Response().SetError(err).Send(ctx)
	})
}
