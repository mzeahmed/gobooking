package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	repo "github.com/mzeahmed/gobooking/internal/adapters/postgresql/sqlc"
)

// uniqueViolation is the PostgreSQL error code raised when a UNIQUE
// constraint (here, users.email) is violated.
const uniqueViolation = "23505"

// Service contains the business logic of the user module. It consumes the
// sqlc-generated Queries directly, with no repository layer in between.
type Service struct {
	pool    *pgxpool.Pool
	queries *repo.Queries
}

// NewService creates a new user service backed by the given pool.
func NewService(pool *pgxpool.Pool) *Service {

	return &Service{
		pool:    pool,
		queries: repo.New(pool),
	}
}

// Create inserts a new user with the default "user" role and returns it
// with its generated fields (ID, CreatedAt, UpdatedAt, Roles) populated.
func (s *Service) Create(ctx context.Context, u User) (User, error) {

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	created, err := q.CreateUser(ctx, repo.CreateUserParams{
		Email:      u.Email,
		Password:   u.PasswordHash,
		FirstName:  pgtype.Text{String: u.FirstName, Valid: u.FirstName != ""},
		LastName:   pgtype.Text{String: u.LastName, Valid: u.LastName != ""},
		IsVerified: u.IsVerified,
	})

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return User{}, ErrEmailTaken
		}

		return User{}, err
	}

	if err := q.AssignDefaultRole(ctx, repo.AssignDefaultRoleParams{
		UserID: created.ID,
		Name:   string(RoleUser),
	}); err != nil {
		return User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return User{}, err
	}

	u.ID = int(created.ID)
	u.CreatedAt = created.CreatedAt.Time
	u.UpdatedAt = created.UpdatedAt.Time
	u.Roles = []Role{RoleUser}

	return u, nil
}

// FindByEmail returns the user with the given email address, or
// ErrNotFound if none exists.
func (s *Service) FindByEmail(ctx context.Context, email string) (User, error) {

	row, err := s.queries.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}

		return User{}, err
	}

	return User{
		ID:           int(row.ID),
		Email:        row.Email,
		PasswordHash: row.Password,
		FirstName:    row.FirstName,
		LastName:     row.LastName,
		IsVerified:   row.IsVerified,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
		Roles:        rolesFromNames(row.Roles),
	}, nil
}

// ListUsers returns every registered user, ordered by ID.
func (s *Service) ListUsers(ctx context.Context) ([]User, error) {

	rows, err := s.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	var users []User

	for _, row := range rows {
		users = append(users, User{
			ID:           int(row.ID),
			Email:        row.Email,
			PasswordHash: row.Password,
			FirstName:    row.FirstName,
			LastName:     row.LastName,
			IsVerified:   row.IsVerified,
			CreatedAt:    row.CreatedAt.Time,
			UpdatedAt:    row.UpdatedAt.Time,
			Roles:        rolesFromNames(row.Roles),
		})
	}

	return users, nil
}

// DeleteUser removes the user identified by req, returning ErrNotFound
// if no user matches the given ID.
func (s *Service) DeleteUser(ctx context.Context, req DeleteRequest) error {

	affected, err := s.queries.DeleteUser(ctx, int32(req.ID))
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrNotFound
	}

	return nil
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
