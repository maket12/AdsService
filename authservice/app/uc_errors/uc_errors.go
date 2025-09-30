package uc_errors

import "errors"

var (
	ErrGetUser           = errors.New("failed to get user")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredential = errors.New("invalid credentials")
	ErrTokenIssue        = errors.New("failed to issue token")
	ErrTokenParse        = errors.New("failed to parse token")
	ErrSessionSave       = errors.New("failed to save session")
	ErrAddUser           = errors.New("failed to add user")
	ErrAddProfile        = errors.New("failed to add profile")
)
