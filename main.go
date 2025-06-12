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
	
	
	db, disconnect, client := data.ConnectDB()
	defer disconnect()
	
	store := data.NewMongoStore(db)
	
	h := handler.NewHandler(db, store, client)

	m := middleware.NewMiddleware(store)
	
	mux := new(middleware.CustomMux)
	
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"OPTIONS","GET","POST","DELETE","PATCH"},
		AllowedHeaders: []string{"Content-Type"},
		AllowCredentials: true,
	})

	mux.RegisterMiddleware(c.Handler)
	mux.RegisterMiddleware(context.ClearHandler)
	
	
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
