package main

import (
	"log"
	"main/common/constan"
	"main/common/mongodb"
	greeting "main/greeting/route"
	invitation "main/invitation/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * constan.MBSize,
	})
	app.Use(cors.New(cors.Config{
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))
	app.Use(logger.New())
	app.Static("/assets", "./assets")

	err := mongodb.InitConnection()
	if err != nil {
		log.Fatal(err)
		return
	}

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
