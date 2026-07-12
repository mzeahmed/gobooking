package user

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Module wires together the user module's dependencies and exposes its
// HTTP routes.
type Module struct {
	handler *Handler
}

// New builds a user Module with its repository, service and handler
// dependencies initialized.
func New(pool *pgxpool.Pool) *Module {

	repo := NewRepository(pool)
	service := NewService(repo)
	handler := NewHandler(service)

	return &Module{
		handler: handler,
	}
}

// RegisterRoutes registers the user module's routes on the given mux.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("DELETE /users", m.handler.Delete)
}
