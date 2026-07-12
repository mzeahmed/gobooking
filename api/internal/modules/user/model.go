// Package user owns the user domain: the User model and its persistence.
// It is consumed by the auth module for registration and login, and
// exposes its own HTTP routes for user account management (e.g.
// deletion).
package user

import "time"

// Role identifies a permission level, matching the roles.name values
// seeded in the roles table.
type Role string

const (
	RoleAdmin     Role = "admin"
	RoleUser      Role = "user"
	RoleModerator Role = "moderator"
	RoleGuest     Role = "guest"
)

// User represents a row of the users table, along with the roles
// assigned to it through the user_roles table.
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
