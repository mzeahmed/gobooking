// Package router assembles the application's HTTP handler by wiring up the
// routes exposed by each module.
package router

import (
	"net/http"

	"github.com/mzeahmed/gobooking/internal/modules/health"
)

// New builds and returns the application's top-level http.Handler, with all
// module routes registered on a fresh http.ServeMux.
func New() http.Handler {

	mux := http.NewServeMux()

	health.New().RegisterRoutes(mux)

	return mux
}
