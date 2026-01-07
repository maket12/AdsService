package dto

type Logout struct {
	RefreshToken string
}

type LogoutResponse struct {
	Logout bool
}
