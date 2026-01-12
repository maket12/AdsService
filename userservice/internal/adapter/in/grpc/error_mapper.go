package grpc

import (
	"ads/pkg/errs"
	"ads/userservice/internal/app/uc_errors"
	"errors"

	"google.golang.org/grpc/codes"
)

type OutError struct {
	Code    codes.Code
	Message string
	Reason  error
}

func NewOutError(code codes.Code, msg string, reason error) *OutError {
	return &OutError{
		Code:    code,
		Message: msg,
		Reason:  reason,
	}
}

func gRPCError(err error) *OutError {
	var w *uc_errors.WrappedError
	if errors.As(err, &w) {
		switch {
		case errors.Is(w.Public, uc_errors.ErrCreateProfileDB),
			errors.Is(w.Public, uc_errors.ErrGetProfileDB),
			errors.Is(w.Public, uc_errors.ErrUpdateProfileDB):
			return NewOutError(codes.Internal, w.Public.Error(), w.Reason)

		default:
			return NewOutError(codes.Internal, "internal error", w.Reason)
		}
	}

	switch {
	case errors.Is(err, uc_errors.ErrInvalidAccountID):
		return NewOutError(codes.NotFound, err.Error(), nil)

	case errors.Is(err, uc_errors.ErrInvalidProfileData),
		errors.Is(err, uc_errors.ErrInvalidPhoneNumber):
		return NewOutError(codes.InvalidArgument, err.Error(), nil)

	case errors.Is(err, errs.ErrNotAuthenticated):
		return NewOutError(codes.Unauthenticated, err.Error(), nil)
	}

	return NewOutError(codes.Internal, "internal error", nil)
}
