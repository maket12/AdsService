package dto

type ValidateToken struct {
	AccessToken string
}

type ValidateTokenResponse struct {
	UserID uint64
	Valid  bool
}
