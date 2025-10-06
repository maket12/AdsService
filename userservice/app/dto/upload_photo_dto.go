package dto

type UploadPhotoDTO struct {
	UserID      uint64
	Data        []byte
	Filename    string
	ContentType string
}
