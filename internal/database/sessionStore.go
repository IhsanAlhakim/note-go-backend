package database

import (
	"backend/internal/config"
	"net/http"

	"github.com/laziness-coders/mongostore"
	"go.mongodb.org/mongo-driver/mongo"
)

// Cookie Based Authentication
func NewSessionStore(db *mongo.Database, cfg *config.Config) *mongostore.MongoStore {

	maxAge := 86400 * 7 // 1 week
	ensureTTL := true
	authKey := []byte(cfg.SessionAuthKey)
	encryptionKey := []byte(cfg.SessionEncryptionKey)

	store := mongostore.NewMongoStore(db.Collection("sessions"), maxAge, ensureTTL, authKey, encryptionKey)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = true
	store.Options.SameSite = http.SameSiteNoneMode
	return store
}
