package grpc

import (
	"ads/authservice/internal/app/uc_errors"
	"errors"

	"google.golang.org/grpc/codes"
)

func gRPCError(err error) (codes.Code, string, error) {
	var w *uc_errors.WrappedError
	if errors.As(err, &w) {
		switch {
		case errors.Is(w.Public, uc_errors.ErrHashPassword),
			errors.Is(w.Public, uc_errors.ErrCreateAccountDB),
			errors.Is(w.Public, uc_errors.ErrGetAccountByEmailDB),
			errors.Is(w.Public, uc_errors.ErrGetAccountByIDDB),
			errors.Is(w.Public, uc_errors.ErrUpdateAccountDB),
			errors.Is(w.Public, uc_errors.ErrGetAccountRoleDB),
			errors.Is(w.Public, uc_errors.ErrUpdateAccountRoleDB),
			errors.Is(w.Public, uc_errors.ErrCreateRefreshSessionDB),
			errors.Is(w.Public, uc_errors.ErrGetRefreshSessionByIDDB),
			errors.Is(w.Public, uc_errors.ErrRevokeRefreshSessionDB),
			errors.Is(w.Public, uc_errors.ErrCreateAccountRoleDB),
			errors.Is(w.Public, uc_errors.ErrGenerateAccessToken),
			errors.Is(w.Public, uc_errors.ErrGenerateRefreshToken):
			return codes.Internal, w.Public.Error(), w.Reason

		case errors.Is(w.Public, uc_errors.ErrInvalidInput):
			return codes.InvalidArgument, w.Public.Error(), w.Reason

		default:
			return codes.Internal, "internal error", w.Reason
		}
	}

	switch {
	case errors.Is(err, uc_errors.ErrInvalidCredentials),
		errors.Is(err, uc_errors.ErrInvalidAccountID):
		return codes.NotFound, err.Error(), nil

	case errors.Is(err, uc_errors.ErrAccountAlreadyExists):
		return codes.AlreadyExists, err.Error(), nil

	case errors.Is(err, uc_errors.ErrCannotLogin),
		errors.Is(err, uc_errors.ErrCannotAssign),
		errors.Is(err, uc_errors.ErrCannotRevoke):
		return codes.FailedPrecondition, err.Error(), nil

	case errors.Is(err, uc_errors.ErrInvalidAccessToken),
		errors.Is(err, uc_errors.ErrInvalidRefreshToken):
		return codes.Unauthenticated, err.Error(), nil
	}

	return codes.Internal, "internal error", nil
}
