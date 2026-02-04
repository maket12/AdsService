package uc_errors

import "errors"

/*
================ Validation failures ================
*/
var (
	ErrInvalidInput  = errors.New("one or several specified parameters are invalid")
	ErrInvalidAdID   = errors.New("ad id is invalid or ad with this id not found")
	ErrCannotPublish = errors.New("ad has been already published or not available")
	ErrCannotReject  = errors.New("ad has been already published or not available")
	ErrCannotDelete  = errors.New("ad has been already deleted or not published yet")
)

/*
================ MongoDB failures ================
*/
var (
	ErrSaveImagesDB = errors.New("failed to save images using db")
	ErrGetImagesDB  = errors.New("failed to get images using db")
)

/*
================ Postgres failures ================
*/
var (
	ErrCreateAdDB       = errors.New("failed to create ad using db")
	ErrGetAdDB          = errors.New("failed to get add using db")
	ErrUpdateAdDB       = errors.New("failed to update ad using db")
	ErrUpdateAdStatusDB = errors.New("failed to update ad status using db")
	ErrDeleteAllAdsDB   = errors.New("failed to all ads using db")
)
