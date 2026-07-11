// Package main is the entry point of the gobooking HTTP API server.
//
// It loads the application configuration, builds the HTTP router and
// starts listening for incoming requests.
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/mzeahmed/gobooking/internal/config"
	"github.com/mzeahmed/gobooking/internal/db"
	"github.com/mzeahmed/gobooking/internal/router"
)

// main loads the configuration, wires up the router and starts the HTTP
// server. It terminates the process if the server fails to start.
func main() {

	cfg := config.Load()

	pool, err := db.New(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	r := router.New(pool, cfg.JWTSecret)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	log.Printf("Server listening on :%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
