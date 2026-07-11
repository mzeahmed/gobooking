// Package auth handles user registration and login, issuing JWT access
// tokens on success. It builds on top of the user module's repository
// rather than owning user persistence itself.
package auth

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mzeahmed/gobooking/internal/modules/user"
)

// Module wires together the auth module's dependencies and exposes its
// HTTP routes.
type Module struct {
	handler *Handler
}

// New builds an auth Module with its repository, service and handler
// dependencies initialized.
func New(pool *pgxpool.Pool, jwtSecret string) *Module {

	repo := user.NewRepository(pool)
	service := NewService(repo, jwtSecret)
	handler := NewHandler(service)

	return &Module{
		handler: handler,
	}
}

// RegisterRoutes registers the auth module's routes on the given mux.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("POST /auth/register", m.handler.Register)
	mux.HandleFunc("POST /auth/login", m.handler.Login)
}
