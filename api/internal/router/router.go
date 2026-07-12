// Package router assembles the application's HTTP handler by wiring up the
// routes exposed by each module.
package router

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mzeahmed/gobooking/internal/modules/auth"
	"github.com/mzeahmed/gobooking/internal/modules/health"
	"github.com/mzeahmed/gobooking/internal/modules/user"
)

// New builds and returns the application's top-level http.Handler, with all
// module routes registered on a fresh http.ServeMux.
func New(pool *pgxpool.Pool, jwtSecret string) http.Handler {

	mux := http.NewServeMux()

	health.New(jwtSecret).RegisterRoutes(mux)
	auth.New(pool, jwtSecret).RegisterRoutes(mux)
	user.New(pool).RegisterRoutes(mux)

	return mux
}
