package router

import (
	"net/http"

	"github.com/mzeahmed/go-booking/internal/health"
	"github.com/mzeahmed/go-booking/internal/middleware"
)

func New() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health.Handler)

	var handler http.Handler = mux

	handler = middleware.Logging(handler)
	handler = middleware.Recovery(handler)

	return handler
}
