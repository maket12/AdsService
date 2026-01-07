package dto

type Login struct {
	Email     string
	Password  string
	IP        *string
	UserAgent *string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}
