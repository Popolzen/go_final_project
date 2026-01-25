package models

import "errors"

var (
	ErrInvalidUsername     = errors.New("username must be at least 3 characters")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserExists          = errors.New("user already exists")
	ErrInvalidSecretType   = errors.New("invalid secret type")
	ErrEmptySecretName     = errors.New("secret name cannot be empty")
	ErrEmptySecretData     = errors.New("secret data cannot be empty")
	ErrSecretNotFound      = errors.New("secret not found")
	ErrUnauthorized        = errors.New("unauthorized access")
	ErrInvalidPassword     = errors.New("password must be at least 8 characters")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidSecretData   = errors.New("invalid secret data")
	ErrSecretAlreadyExists = errors.New("secret already exists")
	ErrInvalidToken        = errors.New("invalid token")
)
