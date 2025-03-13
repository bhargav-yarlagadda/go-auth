package handlers

import (
	"context"
	"server/config"
	"server/db"
	"server/models"
	"strings"
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



// SignIn handler to authenticate user and return JWT
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

	err := userCollection.FindOne(context.TODO(), bson.M{"email": credentials.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email "})
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
	}

	// Generate JWT token
	token, err := config.GenerateSecret(user.ID.Hex(), user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	// Return token 
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"user":user,
	})
}


// Validate Token Route Handler
func Validate(c *fiber.Ctx) error {
	// Get the Authorization header
	tokenHeader := c.Get("Authorization")
	if tokenHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	// Ensure the token follows "Bearer <token>" format
	if !strings.HasPrefix(tokenHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
	}

	// Extract the token value
	tokenString := strings.TrimPrefix(tokenHeader, "Bearer ")

	// Validate the token
	_, err := config.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Token is valid
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Token validated"})
}
// we can send jwt from client to server if we wantt to have authorized routes instead of username and password we use JWT 
// we pass the req to a middleware if its not authorized we cannot move to next routes 

// Middleware to protect routes with JWT
// func AuthMiddleware(c *fiber.Ctx) error {
// 	tokenString := c.Get("Authorization")
// 	if tokenString == "" {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
// 	}

// 	// Remove "Bearer " prefix if present
// 	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

// 	// Validate token
// 	token, err := config.ValidateToken(tokenString)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
// 	}

// 	// Set user info in context
// 	c.Locals("user", token.Claims)

// 	return c.Next() // Allow access
// }
// app.Get("/profile", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
// 	user := c.Locals("user")
// 	return c.JSON(fiber.Map{
// 		"message": "Welcome to your profile",
// 		"user":    user,
// 	})
// })
