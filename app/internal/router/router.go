package router

import (
	"net/http"

	"github.com/mzeahmed/gobooking/internal/modules/health"
)

func New() http.Handler {

	mux := http.NewServeMux()

	health.New().RegisterRoutes(mux)

	return mux
}
