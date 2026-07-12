package user

import "errors"

// DeleteRequest is the expected JSON body of a user deletion request.
type DeleteRequest struct {
	ID int `json:"id"`
}

// Validate checks that the delete request identifies a user by ID.
func (r DeleteRequest) Validate() error {

	if r.ID <= 0 {
		return errors.New("a valid ID is required")
	}

	return nil
}
