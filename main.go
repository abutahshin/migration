package main

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/net/context"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to MongoDB")

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
		log.Printf("Disconnected from MongoDB")
	}()

	coll := client.Database("gen-ai-dev").Collection("access_control")
	userID := "9809b48e-c691-4716-826b-1adf73304e15"

	var result bson.M
	err = coll.FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the user_id %s\n", userID)
		return
	}
	if err != nil {
		log.Fatalf("Error finding document: %v", err)
	}

	// Convert the result to JSON
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling result to JSON: %v", err)
	}
	fmt.Printf("%s\n", jsonData)
}
