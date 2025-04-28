package handler

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ctx = context.TODO()

type Handler struct {
	db *mongo.Database
}

func NewHandler(db *mongo.Database) *Handler {
	return &Handler{db: db}
}
