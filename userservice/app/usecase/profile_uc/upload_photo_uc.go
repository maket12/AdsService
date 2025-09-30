package profile_uc

import (
	"AdsService/userservice/app/dto/profile_dto"
	"AdsService/userservice/app/mappers"
	"AdsService/userservice/app/uc_errors"
	"AdsService/userservice/domain/port"
	"bytes"
)

type UploadPhotoUC struct {
	Profiles port.ProfileRepository
	Photos   port.PhotoRepository
}

func (uc *UploadPhotoUC) Execute(in profile_dto.UploadPhotoDTO) (profile_dto.ProfileResponseDTO, error) {
	if len(in.Data) == 0 {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrEmptyDataPhoto
	}
	if in.Filename == "" {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrEmptyFilenamePhoto
	}
	if in.ContentType == "" {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrEmptyContentTypePhoto
	}

	reader := bytes.NewReader(in.Data)
	objectHexID, err := uc.Photos.UploadPhoto(in.UserID, in.Filename, in.ContentType, reader, int64(len(in.Data)))
	if err != nil || objectHexID == "" {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrMongoUploadPhoto
	}

	profile, err := uc.Profiles.UpdateProfilePhoto(in.UserID, objectHexID)
	if err != nil {
		return profile_dto.ProfileResponseDTO{}, uc_errors.ErrUpdatePhoto
	}

	return mappers.MapIntoProfileDTO(profile), nil
}
