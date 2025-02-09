package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectToDb()  {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env.local file")
	}

	// Print the MONGO_URI to check if it's loaded
	mongoURI := os.Getenv("MONGO_URI")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal("MongoDB Connection Error:", err)
	}

	// Check if the connection is successful
	err = client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err()
	if err != nil {
		log.Fatal("Ping Error:", err)
	}

	fmt.Println("Connected to MongoDB successfully!")
	MongoClient = client 
}
