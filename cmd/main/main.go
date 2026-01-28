package main

import (
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	middleware "backend/internal/middlewares"
	"backend/internal/mux"
	"backend/internal/routes"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	go func() {
		log.Println("Server started at localhost:" + cfg.Port)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
		log.Println("Stopped serving new connection.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete")
}
