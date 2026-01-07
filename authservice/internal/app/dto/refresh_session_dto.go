package dto

type RefreshSession struct {
	OldRefreshToken string
	IP              *string
	UserAgent       *string
}

type RefreshSessionResponse struct {
	AccessToken  string
	RefreshToken string
}
