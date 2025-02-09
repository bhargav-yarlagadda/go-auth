package main

import (
	"context"
	"server/db"
	"server/handlers"

	"github.com/gofiber/fiber/v2"
	
)

func main() {
	app :=fiber.New()
	db.ConnectToDb() 
	app.Post("/auth/sign-up",handlers.SignUp)
	app.Post("/auth/sign-in",handlers.SignIn)

	app.Listen(":8080")


	defer func (){
		db.MongoClient.Disconnect(context.TODO())
	}()
}