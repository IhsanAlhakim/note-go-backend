package main

import (
	"backend/data"
	"backend/middleware"
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/handler"

	"github.com/gorilla/context"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load environment configuration file: %v ", err.Error())
	}
	
	PORT := os.Getenv("PORT")
	
	if PORT == "" {
		PORT = "9000"
	}
	
	
	db, disconnect := data.ConnectDB()
	defer disconnect()
	
	store := data.NewMongoStore(db)
	
	h := handler.NewHandler(db, store)
	
	mux := new(middleware.CustomMux)
	
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"OPTIONS","GET","POST","DELETE","PATCH"},
		AllowedHeaders: []string{"Content-Type"},
		Debug: true,
	})

	mux.RegisterMiddleware(c.Handler)
	mux.RegisterMiddleware(context.ClearHandler)
	
	
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Connection Ok"))
	})

	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/auth", h.GetAuthenticatedUser)
	mux.HandleFunc("/user", h.FindUserById)
	mux.HandleFunc("/create/user", h.CreateUser)
	mux.HandleFunc("/delete/user", h.DeleteUser)
	mux.HandleFunc("/notes", h.FindUserNotes)
	mux.HandleFunc("/note", h.FindNoteById)
	mux.HandleFunc("/create/note", h.CreateNote)
	mux.HandleFunc("/delete/note", h.DeleteNote)
	mux.HandleFunc("/update/note", h.UpdateNote)

	server := new(http.Server)
	server.Addr = ":" + PORT
	server.Handler = mux

	fmt.Println("Server started at localhost:" + PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Println("Shutting down server...")
	}
}
