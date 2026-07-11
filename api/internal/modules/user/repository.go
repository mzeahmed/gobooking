package user

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// uniqueViolation is the PostgreSQL error code raised when a UNIQUE
// constraint (here, users.email) is violated.
const uniqueViolation = "23505"

// Repository provides access to users stored in PostgreSQL.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a user Repository backed by the given pool.
func NewRepository(pool *pgxpool.Pool) *Repository {

	return &Repository{
		pool: pool,
	}
}

// Create inserts a new user and returns it with its generated fields
// (ID, CreatedAt, UpdatedAt) populated.
func (r *Repository) Create(ctx context.Context, u User) (User, error) {

	rolesJSON, err := json.Marshal(u.Roles)
	if err != nil {
		return User{}, err
	}

	query := `
		INSERT INTO users (email, roles, password, first_name, last_name, is_verified)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err = r.pool.QueryRow(
		ctx, query,
		u.Email, rolesJSON, u.PasswordHash, u.FirstName, u.LastName, u.IsVerified,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return User{}, ErrEmailTaken
		}

		return User{}, err
	}

	return u, nil
}

// FindByEmail returns the user with the given email address, or
// ErrNotFound if none exists.
func (r *Repository) FindByEmail(ctx context.Context, email string) (User, error) {

	query := `
		SELECT id, email, roles, password, COALESCE(first_name, ''), COALESCE(last_name, ''), is_verified, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var u User
	var rolesJSON []byte

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &rolesJSON, &u.PasswordHash, &u.FirstName, &u.LastName, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}

		return User{}, err
	}

	if err := json.Unmarshal(rolesJSON, &u.Roles); err != nil {
		return User{}, err
	}

	return u, nil
}
