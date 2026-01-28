package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SessionID             string
	MongoConnectionString string
	SessionAuthKey        string
	SessionEncryptionKey  string
	Port                  string
	AllowedOrigins        string
}

func Load() *Config {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, continue using system environment variables")
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		SessionID:             os.Getenv("SESSION_ID"),
		MongoConnectionString: os.Getenv("MONGO_CONNECTION_STRING"),
		SessionAuthKey:        os.Getenv("SESSION_AUTH_KEY"),
		SessionEncryptionKey:  os.Getenv("SESSION_ENCRYPTION_KEY"),
		Port:                  port,
		AllowedOrigins:        os.Getenv("CLIENT_URL"),
	}
}
