package user

import "errors"

// ErrEmailTaken is returned when the email address is already registered.
var ErrEmailTaken = errors.New("email already registered")

// ErrNotFound is returned when no user matches the requested lookup.
var ErrNotFound = errors.New("user not found")

var ErrFirstNameEmpty = errors.New("Firstname is required")
var ErrLastNameEmpty = errors.New("Lastname is required")
