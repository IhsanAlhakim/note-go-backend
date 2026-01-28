package data

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Database, *mongo.Client) {

	MONGO_CONNECTION_STRING := os.Getenv("MONGO_CONNECTION_STRING")
	if MONGO_CONNECTION_STRING == "" {
		log.Fatal("Database connection string is missing or empty")
	}

	var err error

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGO_CONNECTION_STRING))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v ", err.Error())
	}

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "username", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	db := client.Database("note_go")

	_, err = db.Collection("users").Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}

	return db, client
}
