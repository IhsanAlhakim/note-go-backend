package handler

import (
	"backend/utils"
	"context"

	"github.com/laziness-coders/mongostore"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.TODO()

type R = utils.Response

type Handler struct {
	db    *mongo.Database
	store *mongostore.MongoStore
	client *mongo.Client
}

func NewHandler(db *mongo.Database, store *mongostore.MongoStore, client *mongo.Client) *Handler {
	return &Handler{db: db, store: store, client: client}
}
