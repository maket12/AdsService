package uc_errors

import "errors"

var (
	ErrEnhanceUser           = errors.New("failed to enhance user")
	ErrGetUser               = errors.New("failed to get user")
	ErrUserNotFound          = errors.New("user not found")
	ErrAddProfile            = errors.New("failed to add profile")
	ErrGetProfile            = errors.New("failed to get profile")
	ErrGetProfiles           = errors.New("failed to get profiles")
	ErrProfileNotFound       = errors.New("profile not found")
	ErrUpdateProfile         = errors.New("failed to update profile")
	ErrEmptyDataPhoto        = errors.New("empty file data")
	ErrEmptyFilenamePhoto    = errors.New("filename required")
	ErrEmptyContentTypePhoto = errors.New("content-type required")
	ErrMongoUploadPhoto      = errors.New("failed to upload photo in MongoDB")
	ErrUpdatePhoto           = errors.New("failed to update photo")
	ErrChangeSettings        = errors.New("failed to change settings")
	ErrChangeSubscriptions   = errors.New("failed to change subscriptions")
	ErrGetUserRole           = errors.New("failed to get user role")
	ErrBanUser               = errors.New("failed to ban user")
	ErrUnbanUser             = errors.New("failed to unban user")
	ErrNotAdmin              = errors.New("permission denied")
)
