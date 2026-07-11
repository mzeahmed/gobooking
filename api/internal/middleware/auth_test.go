package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mzeahmed/gobooking/internal/modules/auth"
	"github.com/mzeahmed/gobooking/internal/modules/user"
)

const testSecret = "test-secret"

func signToken(t *testing.T, secret string, userID int, roles []user.Role, expiresAt time.Time) string {
	t.Helper()

	claims := auth.Claims{
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}

	return token
}

func TestAuthenticate_ValidToken(t *testing.T) {

	token := signToken(t, testSecret, 42, []user.Role{user.RoleUser}, time.Now().Add(time.Hour))

	var got AuthUser
	var gotOK bool

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, gotOK = UserFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	Authenticate(testSecret)(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !gotOK {
		t.Fatal("expected AuthUser in context")
	}

	if got.ID != 42 {
		t.Errorf("expected user ID 42, got %d", got.ID)
	}

	if len(got.Roles) != 1 || got.Roles[0] != string(user.RoleUser) {
		t.Errorf("expected roles [%s], got %v", user.RoleUser, got.Roles)
	}
}

func TestAuthenticate_MissingHeader(t *testing.T) {

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	Authenticate(testSecret)(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestAuthenticate_InvalidSignature(t *testing.T) {

	token := signToken(t, "wrong-secret", 42, []user.Role{user.RoleUser}, time.Now().Add(time.Hour))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	Authenticate(testSecret)(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestAuthenticate_ExpiredToken(t *testing.T) {

	token := signToken(t, testSecret, 42, []user.Role{user.RoleUser}, time.Now().Add(-time.Hour))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	Authenticate(testSecret)(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}
