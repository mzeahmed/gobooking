// Package reqctx carries the identity of an authenticated caller through a
// request's context.Context.
//
// It deliberately has no dependency on internal/modules/auth or
// internal/modules/user, so it can be imported both by internal/middleware
// (which depends on internal/modules/auth to validate tokens) and by
// internal/modules/user (which needs the caller's identity to guard
// actions like account deletion) without creating an import cycle:
// middleware -> auth -> user -> middleware.
package reqctx

import "context"

// contextKey is an unexported type to avoid collisions with context keys
// defined in other packages.
type contextKey int

const authUserKey contextKey = iota

// AuthUser is the identity extracted from a validated access token.
type AuthUser struct {
	ID    int
	Roles []string
}

// WithAuthUser returns a copy of ctx carrying u, retrievable later with
// AuthUserFromContext.
func WithAuthUser(ctx context.Context, u AuthUser) context.Context {
	return context.WithValue(ctx, authUserKey, u)
}

// AuthUserFromContext returns the AuthUser attached with WithAuthUser, if
// any.
func AuthUserFromContext(ctx context.Context) (AuthUser, bool) {
	u, ok := ctx.Value(authUserKey).(AuthUser)

	return u, ok
}
