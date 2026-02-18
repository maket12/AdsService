package grpc

import (
	"ads/adservice/internal/app/uc_errors"
	"errors"

	"google.golang.org/grpc/codes"
)

func gRPCError(err error) (codes.Code, string, error) {
	var w *uc_errors.WrappedError
	if errors.As(err, &w) {
		switch {
		case errors.Is(w.Public, uc_errors.ErrSaveImagesDB),
			errors.Is(w.Public, uc_errors.ErrGetImagesDB),
			errors.Is(w.Public, uc_errors.ErrDeleteImagesDB),
			errors.Is(w.Public, uc_errors.ErrCreateAdDB),
			errors.Is(w.Public, uc_errors.ErrGetAdDB),
			errors.Is(w.Public, uc_errors.ErrUpdateAdDB),
			errors.Is(w.Public, uc_errors.ErrUpdateAdStatusDB),
			errors.Is(w.Public, uc_errors.ErrDeleteAdDB),
			errors.Is(w.Public, uc_errors.ErrDeleteAllAdsDB):
			return codes.Internal, w.Public.Error(), w.Reason

		case errors.Is(w.Public, uc_errors.ErrInvalidInput):
			return codes.InvalidArgument, w.Public.Error(), w.Reason

		default:
			return codes.Internal, "internal error", w.Reason
		}
	}

	switch {
	case errors.Is(err, uc_errors.ErrInvalidAdID):
		return codes.NotFound, err.Error(), nil

	case errors.Is(err, uc_errors.ErrCannotPublish),
		errors.Is(err, uc_errors.ErrCannotReject),
		errors.Is(err, uc_errors.ErrCannotDelete):
		return codes.FailedPrecondition, err.Error(), nil
	}

	return codes.Internal, "internal error", nil
}
