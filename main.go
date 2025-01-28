package main

import (
	"fmt"
	"github.com/abutahshin/graphql/db"
	"github.com/abutahshin/graphql/migration"
	"github.com/goccy/go-json"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	client, err := db.ConnectToMongoDB(uri)

	defer db.DisconnectMongoDB(client)

	var dbname string = "dbname"
	var collname string = "collectionname"
	//tahke database name and collection name as input
	//fmt.Println("Enter DB Name: ")
	//fmt.Scanln(&dbname)
	//fmt.Println("Enter Collection Name: ")
	//fmt.Scanln(&collname)

	coll := client.Database(dbname).Collection(collname)

	var userID string

	fmt.Println("Enter Object ID For which you want to Fetch Data: ")
	fmt.Scanln(&userID)
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Fatal(err)
	}
	result, err := db.FetchData(coll, objID)
	if err != nil {
		log.Fatalf("Error finding document: %v", err)
	}

	// Convert the result to JSON
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling result to JSON: %v", err)
	}
	fmt.Printf("%s\n", jsonData)

	//Update time code from here

	// Time offset (6 hours in milliseconds)
	time := 6 * 60 * 60 * 1000

	// Update the date fields
	updatedFields := bson.M{}
	migration.UpdateDate(result, updatedFields, "fieldName", time)
	migration.UpdateDate(result, updatedFields, "fieldName", time)

	if len(updatedFields) > 0 {
		_, err = coll.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": updatedFields})
		if err != nil {
			panic(err)
		}
		fmt.Println("Document updated successfully!")
	} else {
		fmt.Println("No date fields to update.")
	}
}
