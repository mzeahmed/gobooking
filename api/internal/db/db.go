// Package db manages the application's PostgreSQL connection pool.
package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mzeahmed/gobooking/internal/config"
)

// New builds a connection pool from the given configuration.
func New(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
