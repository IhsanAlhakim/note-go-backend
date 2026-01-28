package database

import (
	"backend/internal/config"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(cfg *config.Config) (*mongo.Database, *mongo.Client) {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.MongoConnectionString))
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
