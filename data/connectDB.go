package data

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDB() (*mongo.Database, func() error, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(("Error loading .env file"))
		return nil, nil, err
	}

	MONGO_CONNECTION_STRING := os.Getenv("MONGO_CONNECTION_STRING")
	if MONGO_CONNECTION_STRING == "" {
		log.Fatal("Connection String Cannot Be Empty")
		return nil, nil, err
	}

	client, err := mongo.Connect(options.Client().ApplyURI(MONGO_CONNECTION_STRING))
	if err != nil {
		return nil, nil, err
	}

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "username", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err = client.Database("note_go").Collection("users").Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}

	var disconnectDB = func() error {
		err := client.Disconnect(context.TODO())
		return err
	}

	return client.Database("note_go"), disconnectDB, nil
}
