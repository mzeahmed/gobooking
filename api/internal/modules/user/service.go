package user

import "context"

// Service contains the business logic of the user module.
type Service struct {
	users *Repository
}

// NewService creates a new user service.
func NewService(users *Repository) *Service {

	return &Service{
		users: users,
	}
}

// DeleteUser removes the user identified by req, returning ErrNotFound
// if no user matches the given ID.
func (s *Service) DeleteUser(ctx context.Context, req DeleteRequest) error {

	return s.users.Delete(ctx, req.ID)
}
