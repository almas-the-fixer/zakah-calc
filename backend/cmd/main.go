package main

import (
	"log"

	"github.com/almas-the-fixer/zakah-calc/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	// Swagger Stuff
	_ "github.com/almas-the-fixer/zakah-calc/docs"
	"github.com/gofiber/swagger"
)

// @title           Zakah Calculator API
// @version         1.0
// @description     A clean Zakah calculator API with currency conversion.
// @host            localhost:8080
// @BasePath        /
func main() {

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: No .env file found. Relying on System Environment Variables.")
	} else {
		log.Println("Loaded .env file successfully!")
	}
	app := fiber.New()

	app.Use(cors.New())

	// Swagger Endpoint
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Zakah API is Running!...", "status": "success"})
	})

	// Main Endpoint For Calculating Zakah //
	app.Post("/calculate-zakah", handlers.CalculateZakah)

	log.Fatal(app.Listen(":8080"))
}
