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
//
// authenticate guards routes that require a logged-in caller; the caller
// (see router.New) is expected to pass middleware.Authenticate(jwtSecret).
// It is injected rather than constructed here because
// internal/middleware depends on internal/modules/auth, which itself
// depends on this package (user.Repository, user.User) — importing
// middleware directly from this package would create an import cycle
// (middleware -> auth -> user -> middleware).
func (m *Module) RegisterRoutes(mux *http.ServeMux, authenticate func(http.Handler) http.Handler) {

	// Listing users requires the caller to be authenticated; Handler.List
	// further restricts it to admins, since the role check depends on the
	// caller's own roles, not just the fact that they're logged in.
	mux.Handle(
		"GET /users",
		authenticate(http.HandlerFunc(m.handler.List)),
	)

	// Only a logged-in caller can delete a user account, and only their
	// own account (see Handler.Delete, which reads the target ID from the
	// authenticated identity rather than the request body).
	mux.Handle(
		"DELETE /users/delete",
		authenticate(http.HandlerFunc(m.handler.Delete)),
	)
}
