package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/mzeahmed/gobooking/internal/modules/user"
)

// Service contains the business logic of the auth module.
type Service struct {
	users     *user.Repository
	jwtSecret string
}

// NewService creates a new auth service.
func NewService(users *user.Repository, jwtSecret string) *Service {

	return &Service{
		users:     users,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user with a hashed password and returns an
// access token for it.
func (s *Service) Register(ctx context.Context, req RegisterRequest) (AuthResponse, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{}, err
	}

	created, err := s.users.Create(ctx, user.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Roles:        []user.Role{user.RoleUser},
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsVerified:   false,
	})

	if err != nil {
		return AuthResponse{}, err
	}

	return s.buildAuthResponse(created)
}

// Login verifies the given credentials and returns an access token on
// success.
func (s *Service) Login(ctx context.Context, req LoginRequest) (AuthResponse, error) {

	u, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return AuthResponse{}, ErrInvalidCredentials
		}

		return AuthResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		return AuthResponse{}, ErrInvalidCredentials
	}

	return s.buildAuthResponse(u)
}

// buildAuthResponse issues a token for u and assembles the response
// returned to the client.
func (s *Service) buildAuthResponse(u user.User) (AuthResponse, error) {

	token, err := generateToken(s.jwtSecret, u)
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		Token: token,
		User: UserResponse{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		},
	}, nil
}
