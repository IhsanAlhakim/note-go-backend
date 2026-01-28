package main

import (
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	middleware "backend/internal/middlewares"
	"backend/internal/mux"
	"backend/internal/routes"
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

	routes.Register(mux, m, h)

	server := new(http.Server)
	server.Addr = ":" + cfg.Port
	server.Handler = mux

	fmt.Println("Server started at localhost:" + cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Println("Shutting down server...")
	}
}
