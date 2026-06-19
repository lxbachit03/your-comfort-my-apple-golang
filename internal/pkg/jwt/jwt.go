package jwt

type JwtService interface {
	GenerateAccessToken() (AccessToken, error)
	GenerateRefreshToken() (string, error)
	StoreRefreshToken(refreshToken string) error
}

type AccessToken struct {
	Token string
	TTL   int
}

type jwtService struct{}

// StoreRefreshToken implements [JwtService].
func (jv *jwtService) StoreRefreshToken(refreshToken string) error {
	panic("unimplemented")
}

func NewJwtService() JwtService {
	return &jwtService{}
}

func (jv *jwtService) GenerateAccessToken() (AccessToken, error) {
	panic("unimplemented")
}

func (jv *jwtService) GenerateRefreshToken() (string, error) {
	panic("unimplemented")
}
