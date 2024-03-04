package response

import "github.com/gofiber/fiber/v2"

func Success(ctx *fiber.Ctx, data interface{}) error {

	return ctx.JSON(fiber.Map{
		"data": data,
	})
}
