package main

import (
	"backend/internal/database"
	"backend/internal/handlers"
	middleware "backend/internal/middlewares"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	gctx "github.com/gorilla/context"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, continue using system environment variables")
		}
	}

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "9000"
	}

	db, client := database.Connect()
	defer client.Disconnect(context.TODO())

	store := database.NewSessionStore(db)

	h := handlers.NewHandler(db, store, client)

	m := middleware.New(store)

	mux := new(middleware.CustomMux)

	allowedOrigin := os.Getenv("CLIENT_URL")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{allowedOrigin},
		AllowedMethods:   []string{"OPTIONS", "GET", "POST", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	mux.RegisterMiddleware(c.Handler)
	mux.RegisterMiddleware(gctx.ClearHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Connection Ok"))
	})

	mux.HandleFunc("/login", h.Login)
	mux.Handle("/logout", m.AuthMiddleware(http.HandlerFunc(h.Logout)))
	mux.Handle("/user", m.AuthMiddleware(http.HandlerFunc(h.GetLoggedInUser)))
	mux.HandleFunc("/create/user", h.CreateUser)
	mux.Handle("/delete/user", m.AuthMiddleware(http.HandlerFunc(h.DeleteUser)))
	mux.Handle("/notes", m.AuthMiddleware(http.HandlerFunc(h.FindUserNotes)))
	mux.HandleFunc("/note", h.FindNoteById)
	mux.Handle("/create/note", m.AuthMiddleware(http.HandlerFunc(h.CreateNote)))
	mux.Handle("/delete/note", m.AuthMiddleware(http.HandlerFunc(h.DeleteNote)))
	mux.Handle("/update/note", m.AuthMiddleware(http.HandlerFunc(h.UpdateNote)))

	server := new(http.Server)
	server.Addr = ":" + PORT
	server.Handler = mux

	fmt.Println("Server started at localhost:" + PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Println("Shutting down server...")
	}
}
