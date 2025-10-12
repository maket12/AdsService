package dto

type UploadPhoto struct {
	UserID      uint64
	Data        []byte
	Filename    string
	ContentType string
}
