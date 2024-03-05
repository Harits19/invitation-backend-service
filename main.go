package main

import (
	"main/invitation"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	invitation.Route(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("test")
	})

	app.Listen(":8080")
	defer app.Shutdown()
}
