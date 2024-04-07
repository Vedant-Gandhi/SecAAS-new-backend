package errors

import "errors"

var (
	ErrInvalidEmail = errors.New("invalid user email found")
	ErrUserNotFound = errors.New("user not found")

	ErrUnknown = errors.New("an unknown error has occurred")
)
