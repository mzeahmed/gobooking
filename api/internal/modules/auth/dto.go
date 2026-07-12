package auth

import (
	"errors"
	"strings"

	"github.com/mzeahmed/gobooking/internal/modules/user"
)

// RegisterRequest is the expected JSON body of a registration request.
type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Validate checks that the registration request contains usable data.
func (r RegisterRequest) Validate() error {

	if strings.TrimSpace(r.Email) == "" || !strings.Contains(r.Email, "@") {
		return errors.New("a valid email is required")
	}

	if len(r.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if r.FirstName == "" {
		return user.ErrFirstNameEmpty
	}

	if r.LastName == "" {
		return user.ErrLastNameEmpty
	}

	return nil
}

// LoginRequest is the expected JSON body of a login request.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate checks that the login request contains usable data.
func (r LoginRequest) Validate() error {

	if strings.TrimSpace(r.Email) == "" || strings.TrimSpace(r.Password) == "" {
		return errors.New("email and password are required")
	}

	return nil
}

// UserResponse is the public representation of a user, safe to return
// in HTTP responses.
type UserResponse struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// AuthResponse is returned by both registration and login on success.
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
