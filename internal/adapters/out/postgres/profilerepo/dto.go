package profilerepo

type RegisterDTO struct {
	Email    string
	Password string
}

type RegisterResponseDTO struct {
	AccessToken  string
	RefreshToken string
}

type LoginDTO struct {
	Email    string
	Password string
}

type LoginResponseDTO struct {
	AccessToken  string
	RefreshToken string
}

type ValidateTokenDTO struct {
	AccessToken string
}

type ValidateTokenResponseDTO struct {
	UserID uint64
	Valid  bool
}
