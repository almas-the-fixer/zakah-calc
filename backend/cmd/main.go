package main

import (
	"log"

	"github.com/almas-the-fixer/zakah-calc/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Zakah API is Running!...", "status": "success"})
	})

	// Main Endpoint For Calculating Zakah //
	app.Post("/calculate-zakah", handlers.CalculateZakah)

	log.Fatal(app.Listen(":8080"))
}
