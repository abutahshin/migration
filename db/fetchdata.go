package db

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/net/context"
)

func FetchData(coll *mongo.Collection, objID primitive.ObjectID) (bson.M, error) {
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the ObjectID %s\n", objID)
		return nil, nil
	}
	return result, nil
}
