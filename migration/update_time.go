package migration

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

func MongoDbUpdateDate(result bson.M, updatedFields bson.M, fieldName string, timeOffset int) {
	if date, ok := result[fieldName].(primitive.DateTime); ok {
		// Subtract 6 hours from the date
		updatedDate := time.Unix(0, int64(date)*int64(time.Millisecond)).Add(-time.Duration(timeOffset) * time.Millisecond)
		updatedFields[fieldName] = primitive.NewDateTimeFromTime(updatedDate)
	}
}
