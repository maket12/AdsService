package port

import (
	"context"
	"io"
)

type PhotoRepository interface {
	UploadPhoto(ctx context.Context, userID uint64, title, contentType string, r io.Reader, size int64) (string, error)
}
