package main

import (
	"main/invitation"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	invitation.Route(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("test")
	})

	app.Listen(":8080")
}
