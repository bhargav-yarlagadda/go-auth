package handlers

import (
	"context"
	"server/db"
	"server/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error {
	var user models.User 
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "All fields are required"})
	}

	// Hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}
	user.Password = string(hashedPassword) 

	// Get MongoDB collection
	userCollection := db.MongoClient.Database("app").Collection("users") 

	// Check if user already exists
	var foundUser models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&foundUser)
	if err == nil { // If no error, user already exists
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already in use"})
	}

	// Set auto-generated fields
	user.ID = primitive.NewObjectID()
	user.CreatedTime = time.Now()

	// Insert new user
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

func SignIn(c *fiber.Ctx) error {
	var credentials models.SignInData 
	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if credentials.Email == "" || credentials.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
	}

	userCollection := db.MongoClient.Database("app").Collection("users")
	var user models.User

	// FIXED: Decode into a pointer
	err := userCollection.FindOne(context.TODO(), bson.M{"email": credentials.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// FIXED: Removed space in `[]byte(credentials.Password)`
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login successful"})
}
