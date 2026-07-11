// Package user owns the user domain: the User model and its persistence.
// It is consumed by the auth module rather than exposing its own HTTP
// routes.
package user

import "time"

// Role identifies a permission level, matching the roles used in the
// source Symfony application.
type Role string

const (
	RoleUser    Role = "ROLE_USER"
	RoleManager Role = "ROLE_MANAGER"
	RoleAdmin   Role = "ROLE_ADMIN"
)

// User represents a row of the users table.
type User struct {
	ID           int
	Email        string
	PasswordHash string
	Roles        []Role
	FirstName    string
	LastName     string
	IsVerified   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
