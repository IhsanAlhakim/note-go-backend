package middleware

import (
	"backend/internal/config"

	"github.com/laziness-coders/mongostore"
)

type Middleware struct {
	store *mongostore.MongoStore
	cfg   *config.Config
}

func New(store *mongostore.MongoStore, cfg *config.Config) *Middleware {
	return &Middleware{store: store, cfg: cfg}
}
