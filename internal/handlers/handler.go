package handlers

import (
	"backend/internal/config"
	"context"

	"github.com/laziness-coders/mongostore"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.TODO()

type Handler struct {
	db     *mongo.Database
	store  *mongostore.MongoStore
	client *mongo.Client
	cfg    *config.Config
}

func New(db *mongo.Database, store *mongostore.MongoStore, client *mongo.Client, config *config.Config) *Handler {
	return &Handler{db: db, store: store, client: client, cfg: config}
}
