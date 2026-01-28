package data

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/laziness-coders/mongostore"
	"go.mongodb.org/mongo-driver/mongo"
)

var SESSION_ID, authKey, encryptionKey string

func init() {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, continue using system environment variables")
		}
	}

	SESSION_ID = os.Getenv("SESSION_ID")
	authKey = os.Getenv("SESSION_AUTH_KEY")
	encryptionKey = os.Getenv("SESSION_ENCRYPTION_KEY")

	if SESSION_ID == "" || authKey == "" || encryptionKey == "" {
		log.Fatal("Session credentials is missing")
	}
}

// Cookies created automatically after session is created
// Cookie Based Authentication
func NewMongoStore(db *mongo.Database) *mongostore.MongoStore {
	maxAge := 86400 * 7 // 1 week
	ensureTTL := true
	authKey := []byte(authKey)
	encryptionKey := []byte(encryptionKey)

	store := mongostore.NewMongoStore(db.Collection("sessions"), maxAge, ensureTTL, authKey, encryptionKey)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = true                    // hanya untuk dev
	store.Options.SameSite = http.SameSiteNoneMode // none untuk prod
	return store
}
