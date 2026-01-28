package middleware

import (
	"backend/utils"

	"github.com/laziness-coders/mongostore"
)


type R = utils.Response

type Middleware struct {
	store *mongostore.MongoStore
}

func NewMiddleware(store *mongostore.MongoStore) *Middleware {
	return &Middleware{store: store}
}
