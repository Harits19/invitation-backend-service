package main

import (
	greeting "main/greeting/route"
	invitation "main/invitation/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	invitationRoute := app.Group("/invitation")
	greetingRoute := app.Group("/greeting")

	invitation.Route(invitationRoute)
	greeting.Route(greetingRoute)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("test")
	})

	app.Listen(":8080")
	defer app.Shutdown()
}
