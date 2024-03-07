package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func Success(ctx *fiber.Ctx, data interface{}) error {
	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    data,
	})
}

func Error(ctx *fiber.Ctx, err error, status *int) error {
	var finalStatus = http.StatusInternalServerError
	if status != nil {
		finalStatus = *status
	}
	return ctx.Status(finalStatus).JSON(fiber.Map{
		"error": err.Error(),
	})
}
