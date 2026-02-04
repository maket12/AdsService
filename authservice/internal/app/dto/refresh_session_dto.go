package dto

type RefreshSessionInput struct {
	OldRefreshToken string
	IP              *string
	UserAgent       *string
}

type RefreshSessionOutput struct {
	AccessToken  string
	RefreshToken string
}
