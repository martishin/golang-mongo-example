package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Get MongoDB connection string from environment variable
	mongoURI := os.Getenv("MONGODB_URI")

	// Connect to the MongoDB cluster
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Select the database and collection
	collection := client.Database("tweet_centre").Collection("tweets")

	// Create a new context with a 10-second timeout for operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert some tweet data
	tweets := []interface{}{
		bson.D{{"username", "James"}, {"tweet", "This is my first tweet"}, {"date", "2024-08-01"}},
		bson.D{{"username", "James"}, {"tweet", "Enjoying the sunshine today!"}, {"date", "2024-08-02"}},
		bson.D{{"username", "James"}, {"tweet", "Learning Go is fun!"}, {"date", "2024-08-03"}},
		bson.D{{"username", "Bob"}, {"tweet", "Started a new project today."}, {"date", "2024-08-01"}},
		bson.D{{"username", "Bob"}, {"tweet", "MongoDB is great for storing documents."}, {"date", "2024-08-02"}},
		bson.D{{"username", "Emily"}, {"tweet", "Having a great time coding."}, {"date", "2024-08-01"}},
		bson.D{{"username", "Emily"}, {"tweet", "Exploring the new features of Go."}, {"date", "2024-08-03"}},
	}

	// Insert the tweets into the collection
	_, err = collection.InsertMany(ctx, tweets)
	if err != nil {
		log.Fatalf("Error inserting data: %v", err)
	}

	// Calculate the average tweet length for James
	var result bson.M
	matchStage := bson.D{{"$match", bson.D{{"username", "James"}}}}
	groupStage := bson.D{
		{
			"$group", bson.D{
				{"_id", nil},
				{"average_length", bson.D{{"$avg", bson.D{{"$strLenCP", "$tweet"}}}}},
			},
		},
	}
	pipeline := mongo.Pipeline{matchStage, groupStage}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatalf("Error calculating average tweet length for James: %v", err)
	}
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			log.Fatalf("Error decoding result: %v", err)
		}
		fmt.Printf("Average tweet length for James was: %.2f characters\n", result["average_length"])
	}

	// Calculate the average tweet length for all users
	groupStageAll := bson.D{
		{
			"$group", bson.D{
				{"_id", nil},
				{"average_length", bson.D{{"$avg", bson.D{{"$strLenCP", "$tweet"}}}}},
			},
		},
	}
	cursorAll, err := collection.Aggregate(ctx, mongo.Pipeline{groupStageAll})
	if err != nil {
		log.Fatalf("Error calculating average tweet length for all users: %v", err)
	}
	if cursorAll.Next(ctx) {
		err = cursorAll.Decode(&result)
		if err != nil {
			log.Fatalf("Error decoding result: %v", err)
		}
		fmt.Printf("Average tweet length for all users was: %.2f characters\n", result["average_length"])
	}
}
