package user

import (
	"context"
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

// Create inserts a new user with the default "user" role and returns it
// with its generated fields (ID, CreatedAt, UpdatedAt, Roles) populated.
func (r *Repository) Create(ctx context.Context, u User) (User, error) {

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback(ctx)

	insertUser := `
		INSERT INTO users (email, password, first_name, last_name, is_verified)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(
		ctx,
		insertUser,
		u.Email,
		u.PasswordHash,
		u.FirstName,
		u.LastName,
		u.IsVerified,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return User{}, ErrEmailTaken
		}

		return User{}, err
	}

	assignDefaultRole := `
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, id FROM roles WHERE name = $2
	`

	if _, err := tx.Exec(ctx, assignDefaultRole, u.ID, string(RoleUser)); err != nil {
		return User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return User{}, err
	}

	u.Roles = []Role{RoleUser}

	return u, nil
}

// FindByEmail returns the user with the given email address, or
// ErrNotFound if none exists.
func (r *Repository) FindByEmail(ctx context.Context, email string) (User, error) {

	query := `
		SELECT u.id, u.email, u.password, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''),
		       u.is_verified, u.created_at, u.updated_at,
		       COALESCE(array_agg(r.name) FILTER (WHERE r.name IS NOT NULL), '{}')
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.id
		LEFT JOIN roles r ON r.id = ur.role_id
		WHERE u.email = $1
		GROUP BY u.id
	`

	var u User
	var roleNames []string

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt, &roleNames,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}

		return User{}, err
	}

	u.Roles = rolesFromNames(roleNames)

	return u, nil
}

// List returns every user, ordered by ID.
func (r *Repository) List(ctx context.Context) ([]User, error) {

	query := `
		SELECT u.id, u.email, u.password, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''),
		       u.is_verified, u.created_at, u.updated_at,
		       COALESCE(array_agg(r.name) FILTER (WHERE r.name IS NOT NULL), '{}')
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.id
		LEFT JOIN roles r ON r.id = ur.role_id
		GROUP BY u.id
		ORDER BY u.id
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		var roleNames []string

		if err := rows.Scan(
			&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt, &roleNames,
		); err != nil {
			return nil, err
		}

		u.Roles = rolesFromNames(roleNames)

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// rolesFromNames converts role names as scanned from the roles table into
// the Role type used by the domain model.
func rolesFromNames(names []string) []Role {

	roles := make([]Role, len(names))
	for i, name := range names {
		roles[i] = Role(name)
	}

	return roles
}

// Delete removes the user matching the given ID, returning ErrNotFound
// if no such user exists.
func (r *Repository) Delete(ctx context.Context, id int) error {

	query := `DELETE FROM users WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
