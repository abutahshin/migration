package migration

import (
	"fmt"
	arangodb "github.com/arangodb/go-driver"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"log"
	"time"
)

func MongoDbUpdateDate(ctx context.Context, fieldsToUpdate []string, nofDoc int, db *mongo.Database, collection *mongo.Collection, collectionName string, timeOffset int) {
	// Build the filter conditions for MongoDB query
	filterConditions := bson.M{}
	for _, field := range fieldsToUpdate {
		// MongoDB query condition: field must exist and be non-null
		filterConditions[field] = bson.M{"$ne": nil}
	}

	// Create a cursor for the collection with the filter conditions
	cursor, err := collection.Find(ctx, filterConditions)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	defer cursor.Close(ctx)

	updateCount := 0

	// Iterate over the cursor to process each document
	for cursor.Next(ctx) {
		// Stop after updating `nofDoc` documents
		if updateCount >= nofDoc && nofDoc != -1 {
			log.Printf("Stopped after %d updated\n", nofDoc)
			break
		}

		var doc map[string]interface{}
		err := cursor.Decode(&doc)
		if err != nil {
			log.Printf("Failed to decode document: %v\n", err)
			continue
		}

		// Create an update document for MongoDB
		updateDoc := bson.M{}

		// Process each field to update the date
		for _, field := range fieldsToUpdate {
			if dateStr, ok := doc[field].(string); ok {
				fieldTime, err := time.Parse(time.RFC3339, dateStr)
				if err != nil {
					log.Printf("Error parsing time for document: %v\n", err)
					continue
				}

				// Subtract 6 hours to convert BDT to UTC
				fieldTimeUTC := fieldTime.Add(-time.Duration(timeOffset) * time.Millisecond)

				// Reformat the updated time and assign it to the update document
				updateDoc[field] = fieldTimeUTC.Format(time.RFC3339)
			} else {
				// Handle case where the field does not exist or is not a valid string
				log.Printf("Field %s does not exist or is not a valid string in document\n", field)
			}
		}

		// If there are fields to update, proceed with the update operation
		if len(updateDoc) > 0 {
			// Use the `_id` or any unique identifier for the document to update
			_, err := collection.UpdateOne(
				ctx,
				bson.M{"_id": doc["_id"]}, // Assuming _id is the unique identifier
				bson.M{"$set": updateDoc},
				//options.Update().SetUpsert(false), // Don't create a new document if it doesn't exist
			)
			if err != nil {
				log.Printf("Failed to update document: %v\n", err)
			} else {
				fmt.Printf("Updated document with new times\n")
				updateCount++
			}
		}
	}

	// If there was an error iterating, log it
	if err := cursor.Err(); err != nil {
		log.Printf("Error during cursor iteration: %v\n", err)
	}

	log.Printf("Total %d documents updated:\n", updateCount)
	log.Println("MongoDB Migration completed.\n")
}
func ArangoDbUpdateDate(ctx context.Context, fieldsToUpdate []string, nofDoc int, db arangodb.Database, collection arangodb.Collection, collectionName string, timeOffset int) {
	filterConditions := ""
	for i, field := range fieldsToUpdate {
		if i > 0 {
			filterConditions += " OR "
		}
		filterConditions += fmt.Sprintf("doc.%s != null", field)
	}

	query := fmt.Sprintf("FOR doc IN %s FILTER %s RETURN doc", collectionName, filterConditions)

	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	defer cursor.Close()

	updateCount := 0

	// Iterate over the cursor to process each document
	for cursor.HasMore() {

		//Stop after update nofDoc documenct
		if updateCount >= nofDoc && nofDoc != -1 {
			log.Printf("Stopped after %d updated\n", nofDoc)
			break
		}

		var doc map[string]interface{}
		meta, err := cursor.ReadDocument(ctx, &doc)
		if err != nil {
			log.Printf("Failed to read document: %v\n", err)
			continue
		}

		// Process each field
		for _, field := range fieldsToUpdate {
			// Check if the field exists in the document
			if dateStr, ok := doc[field].(string); ok {
				fieldTime, err := time.Parse(time.RFC3339, dateStr)
				if err != nil {
					log.Printf("Error parsing time for document %s: %v\n", meta.Key, err)
					continue
				}
				// Subtract 6 hours to convert BDT to UTC
				fieldTimeUTC := fieldTime.Add(-time.Duration(timeOffset) * time.Millisecond)

				// Reformat the updated time and assign it back to the document
				doc[field] = fieldTimeUTC.Format(time.RFC3339)
			} else {
				// Handle case where the field does not exist or is not a string
				log.Printf("Field %s does not exist or is not a valid string in document %s\n", field, meta.Key)
			}
		}

		// Update the document in the collection
		_, err = collection.UpdateDocument(ctx, meta.Key, doc)
		if err != nil {
			log.Printf("Failed to update document %s: %v\n", meta.Key, err)
		} else {
			fmt.Printf("Updated document %s with new times\n", meta.Key)
			updateCount++
		}

	}

	log.Printf("Total %d documents updated:\n", updateCount)
	log.Println("arangoDB Migration completed.\n")
}
