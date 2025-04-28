package main

import (
	"backend/data"
	"backend/middleware"
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/handler"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "9000"
	}

	db, disconnect, err := data.ConnectDB()

	if err != nil {
		log.Fatal("MongoDB connection error")
	}

	h := handler.NewHandler(db)

	mux := new(middleware.CustomMux)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
	mux.HandleFunc("/users", h.FindUserById)
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
		defer disconnect()
	}
}
