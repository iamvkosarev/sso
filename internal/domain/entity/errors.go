package entity

import "errors"

var (
	ErrNoMetadata = errors.New("no metadata found in context")
)

// Token errors
var (
	ErrNoAuthHeader      = errors.New("no authorization header provided")
	ErrInvalidAuthHeader = errors.New("invalid authorization header format")
	ErrTokenExpired      = errors.New("token expired")
	ErrTokenIsInvalid    = errors.New("invalid token")
)

// User errors
var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserExists        = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)
