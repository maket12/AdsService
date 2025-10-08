package usecase

import (
	"ads/userservice/app/dto"
	"ads/userservice/app/mappers"
	"ads/userservice/app/uc_errors"
	"ads/userservice/domain/port"
	"bytes"
)

type UploadPhotoUC struct {
	Profiles port.ProfileRepository
	Photos   port.PhotoRepository
}

func (uc *UploadPhotoUC) Execute(in dto.UploadPhoto) (dto.ProfileResponse, error) {
	if len(in.Data) == 0 {
		return dto.ProfileResponse{}, uc_errors.ErrEmptyDataPhoto
	}
	if in.Filename == "" {
		return dto.ProfileResponse{}, uc_errors.ErrEmptyFilenamePhoto
	}
	if in.ContentType == "" {
		return dto.ProfileResponse{}, uc_errors.ErrEmptyContentTypePhoto
	}

	reader := bytes.NewReader(in.Data)
	objectHexID, err := uc.Photos.UploadPhoto(in.UserID, in.Filename, in.ContentType, reader, int64(len(in.Data)))
	if err != nil || objectHexID == "" {
		return dto.ProfileResponse{}, uc_errors.ErrMongoUploadPhoto
	}

	profile, err := uc.Profiles.UpdateProfilePhoto(in.UserID, objectHexID)
	if err != nil {
		return dto.ProfileResponse{}, uc_errors.ErrUpdatePhoto
	}

	return mappers.MapIntoProfileDTO(profile), nil
}
