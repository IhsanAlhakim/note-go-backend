package handler

import (
	"context"

	"github.com/laziness-coders/mongostore"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.TODO()

type Handler struct {
	db    *mongo.Database
	store *mongostore.MongoStore
}

func NewHandler(db *mongo.Database, store *mongostore.MongoStore) *Handler {
	return &Handler{db: db, store: store}
}
