package uc_errors

import "errors"

var (
	ErrEnhanceUser     = errors.New("failed to enhance user")
	ErrGetUser         = errors.New("failed to get user")
	ErrUserNotFound    = errors.New("user not found")
	ErrGetProfile      = errors.New("failed to get profile")
	ErrGetProfiles     = errors.New("failed to get profiles")
	ErrProfileNotFound = errors.New("profile not found")
	ErrGetUserRole     = errors.New("failed to get user role")
	ErrBanUser         = errors.New("failed to ban user")
	ErrUnbanUser       = errors.New("failed to unban user")
	ErrNotAdmin        = errors.New("permission denied")
)
