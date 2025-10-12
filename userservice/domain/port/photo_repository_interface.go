package port

import "io"

type PhotoRepository interface {
	UploadPhoto(userID uint64, title, contentType string, r io.Reader, size int64) (string, error)
}
