package main

import (
	"context"
	"server/db"
	"server/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Enable CORS with specific origins
	
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // Replace with frontend URL
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true, // Allowed only with specific origins
	}))

	db.ConnectToDb()

	app.Post("/auth/sign-up", handlers.SignUp)
	app.Post("/auth/sign-in", handlers.SignIn)
	app.Get("/auth/validate",handlers.Validate)

	// Start the server
	app.Listen(":8080")

	// Ensure MongoDB disconnects properly
	defer func() {
		db.MongoClient.Disconnect(context.TODO())
	}()
}
