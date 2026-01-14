package service

type Tokens struct {
	AccessToken  string
	RefreshToken string
	UserID       string
}

// AuthService
// ========= Auth Contract =============
type AuthService interface {
	Register(countryCode, phone, password, name string) (Tokens, error)
	Login(phone, password, contryCode string) (Tokens, error)
	RefreshAccessToken(refreshToken string) (Tokens, error)
	Logout(id string) error
}
