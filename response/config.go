package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func Success(ctx *fiber.Ctx, data interface{}) error {

	return ctx.JSON(fiber.Map{
		"data": data,
	})
}

func Error(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"error": err.Error(),
	})
}
