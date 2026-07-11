package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/mzeahmed/gobooking/internal/modules/auth"
	"github.com/mzeahmed/gobooking/internal/response"
)

// contextKey is an unexported type to avoid collisions with context keys
// defined in other packages.
type contextKey int

const authUserContextKey contextKey = iota

// AuthUser is the identity extracted from a validated access token.
type AuthUser struct {
	ID    int
	Roles []string
}

// Authenticate returns a middleware that validates the
// "Authorization: Bearer <token>" header on incoming requests. Requests
// without a valid, unexpired token are rejected with 401 before reaching
// next. On success, the authenticated user is attached to the request
// context and can be retrieved with UserFromContext.
func Authenticate(jwtSecret string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			tokenString, ok := bearerToken(r)
			if !ok {
				response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "missing bearer token"})
				return
			}

			claims, err := auth.ParseToken(jwtSecret, tokenString)
			if err != nil {
				response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
				return
			}

			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token subject"})
				return
			}

			roles := make([]string, len(claims.Roles))
			for i, role := range claims.Roles {
				roles[i] = string(role)
			}

			ctx := context.WithValue(r.Context(), authUserContextKey, AuthUser{
				ID:    userID,
				Roles: roles,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// bearerToken extracts the token from an "Authorization: Bearer <token>"
// header.
func bearerToken(r *http.Request) (string, bool) {

	const prefix = "Bearer "

	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", false
	}

	return token, true
}

// UserFromContext returns the AuthUser attached by Authenticate, if any.
func UserFromContext(ctx context.Context) (AuthUser, bool) {

	u, ok := ctx.Value(authUserContextKey).(AuthUser)

	return u, ok
}
