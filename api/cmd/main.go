// Package main is the entry point of the gobooking HTTP API server.
//
// It loads the application configuration, builds the HTTP router and
// starts listening for incoming requests.
package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mzeahmed/gobooking/internal/config"
	"github.com/mzeahmed/gobooking/internal/router"
)

// main loads the configuration, wires up the router and starts the HTTP
// server. It terminates the process if the server fails to start.
func main() {
	ctx := context.Background()

	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	logger.Info("connected to database", "dsn", cfg.DSN)

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
