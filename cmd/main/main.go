package main

import (
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	middleware "backend/internal/middlewares"
	"backend/internal/mux"
	"context"
	"fmt"
	"log"
	"net/http"

	gctx "github.com/gorilla/context"
	"github.com/rs/cors"
)

func main() {
	cfg := config.Load()

	db, client := database.Connect(cfg)
	defer client.Disconnect(context.TODO())

	store := database.NewSessionStore(db, cfg)

	h := handlers.New(db, store, client, cfg)

	m := middleware.New(store, cfg)

	mux := mux.New()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.AllowedOrigins},
		AllowedMethods:   []string{"OPTIONS", "GET", "POST", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	mux.RegisterMiddleware(c.Handler)
	mux.RegisterMiddleware(gctx.ClearHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Connection Ok"))
	})

	mux.Handle("/login", m.CheckContentType(http.HandlerFunc(h.Login)))
	mux.Handle("/logout", m.Auth(http.HandlerFunc(h.Logout)))
	mux.Handle("/user", m.Auth(http.HandlerFunc(h.GetLoggedInUser)))
	mux.Handle("/create/user", m.CheckContentType(http.HandlerFunc(h.CreateUser)))
	mux.Handle("/delete/user", m.Auth(http.HandlerFunc(h.DeleteUser)))
	mux.Handle("/notes", m.Auth(http.HandlerFunc(h.FindUserNotes)))
	mux.Handle("/note", m.Auth(http.HandlerFunc(h.FindNoteById)))
	mux.Handle("/create/note", m.Auth(m.CheckContentType(http.HandlerFunc(h.CreateNote))))
	mux.Handle("/delete/note", m.Auth(http.HandlerFunc(h.DeleteNote)))
	mux.Handle("/update/note", m.Auth(m.CheckContentType(http.HandlerFunc(h.UpdateNote))))

	server := new(http.Server)
	server.Addr = ":" + cfg.Port
	server.Handler = mux

	fmt.Println("Server started at localhost:" + cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Println("Shutting down server...")
	}
}
