package dto

type ValidateTokenDTO struct {
	AccessToken string
}

type ValidateTokenResponse struct {
	UserID uint64
	Valid  bool
}
