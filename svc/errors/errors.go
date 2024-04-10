package errors

import "errors"

var (
	ErrInvalidEmail   = errors.New("invalid user email found")
	ErrUserNotFound   = errors.New("user not found")
	ErrInviteNotFound = errors.New("invite not found")

	ErrUnknown = errors.New("an unknown error has occurred")

	ErrInvalidPassHash      = errors.New("password hash is not valid")
	ErrInvalidAsymmetricKey = errors.New("Asymmetric Key is not valid")
	ErrInvalidSymmetricKey  = errors.New("Symmetric Key is not valid")

	ErrInvalidID = errors.New("ID is not valid")

	ErrInvalidOrganizationID = errors.New("Orgnization ID is not valid")
)
