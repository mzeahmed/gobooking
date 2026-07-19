package user

import "errors"

// DeleteRequest identifies the user account to delete.
//
// Its ID is populated from the authenticated caller's identity (extracted
// from the JWT by middleware.Authenticate), never from client-supplied
// input. This guarantees a user can only ever delete their own account,
// not an arbitrary one picked by ID.
type DeleteRequest struct {
	ID int
}

// Validate checks that the delete request identifies a user by ID. This is
// a defensive check: ID is expected to always be a valid, positive value
// since it comes from a token we issued ourselves, not from user input.
func (r DeleteRequest) Validate() error {

	if r.ID <= 0 {
		return errors.New("a valid ID is required")
	}

	return nil
}

// Response is the public representation of a user, safe to return in
// HTTP responses: notably, it excludes the password hash.
type Response struct {
	ID        int      `json:"id"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Roles     []string `json:"roles"`
}

// newUserResponse converts a User to its public representation.
func newUserResponse(u User) Response {

	roles := make([]string, len(u.Roles))
	for i, role := range u.Roles {
		roles[i] = string(role)
	}

	return Response{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Roles:     roles,
	}
}
