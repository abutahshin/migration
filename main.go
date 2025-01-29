package main

import (
	"context"
	"fmt"
	_ "fmt"
	"github.com/abutahshin/migration/db"
	"github.com/abutahshin/migration/migration"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
)

const (
	time = 6 * 60 * 60 * 1000
)

func main() {
	MongoDB()
	ArangoDB()
}
func MongoDB() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	client, err := db.ConnectToMongoDB(uri)
	defer db.DisconnectMongoDB(client)

	var dbname string = "dbname"
	var collname string = "collectionname"

	coll := client.Database(dbname).Collection(collname)

	//Update time code from here

	// Time offset (6 hours in milliseconds)
	db := client.Database(dbname)
	// Update the date fields
	fieldsToUpdate := []string{"created_at", "otp_sent_time"}
	migration.MongoDbUpdateDate(context.Background(), fieldsToUpdate, 1, db, coll, collname, time)
}
func ArangoDB() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	arango_URI := os.Getenv("ARANGODB_URI")
	arangoDB_uname := os.Getenv("ARANGODB_USERNAME")
	arangoDB_pwd := os.Getenv("ARANGODB_PASSWORD")

	client, err := db.ConnectToArangoDB(arango_URI, arangoDB_uname, arangoDB_pwd)

	if err != nil {
		log.Fatal(err)
	}
	Db_Name := "shikho"
	db, err := client.Database(nil, Db_Name) // Change database name
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	collName := "accounts"
	//Select collection
	coll, err := db.Collection(nil, collName) // Change collection name
	if err != nil {
		log.Fatalf("Failed to open collection: %v", err)
	}
	fieldsToUpdate := []string{"created_at", "otp_sent_time"} // Change the field name
	migration.ArangoDbUpdateDate(context.Background(), fieldsToUpdate, 1, db, coll, collName, time)
}
