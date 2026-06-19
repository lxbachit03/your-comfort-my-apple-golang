package jwt

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/security/encrypt"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/config"
)

type JwtService interface {
	GenerateAccessToken(data JwtPayload) (AccessToken, error)
	GenerateRefreshToken() (string, error)
	StoreRefreshToken(refreshToken string) error
}

type AccessToken struct {
	Token string
	TTL   int
}

type JwtPayload struct {
	UserUUID  string
	UserEmail string
}

var (
	JWT_SECRET            string
	JWT_ENCRYPTION_KEY    string
	JWT_ACCESS_TOKEN_TTL  time.Duration
	JWT_REFRESH_TOKEN_TTL time.Duration
	JWT_ISSUER            string
)

type jwtService struct{}

func NewJwtService() JwtService {

	JWT_SECRET = config.AppConfig.Security.Jwt.Secret
	JWT_ENCRYPTION_KEY = config.AppConfig.Security.Jwt.EncryptionKey
	JWT_ACCESS_TOKEN_TTL = time.Duration(config.AppConfig.Security.Jwt.AccessTokenTtl) * time.Second
	JWT_REFRESH_TOKEN_TTL = time.Duration(config.AppConfig.Security.Jwt.RefreshTokenTtl) * time.Second
	JWT_ISSUER = config.AppConfig.Security.Jwt.Issuer

	return &jwtService{}
}

func (jv *jwtService) GenerateAccessToken(data JwtPayload) (AccessToken, error) {
	payload := JwtPayload{
		UserUUID:  data.UserUUID,
		UserEmail: data.UserEmail,
	}

	rawData, err := json.Marshal(payload)
	if err != nil {
		return AccessToken{}, err
	}

	encryptedData, err := encrypt.EncryptAES(rawData, []byte(JWT_ENCRYPTION_KEY))
	if err != nil {
		return AccessToken{}, err
	}

	claims := jwt.MapClaims{
		"data": encryptedData,
		"jti":  uuid.NewString(),
		"exp":  time.Now().Add(JWT_ACCESS_TOKEN_TTL).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return AccessToken{}, err
	}

	return AccessToken{
		Token: tokenString,
		TTL:   int(JWT_ACCESS_TOKEN_TTL.Seconds()),
	}, nil
}

func (jv *jwtService) GenerateRefreshToken() (string, error) {
	panic("unimplemented")
}

func (jv *jwtService) StoreRefreshToken(refreshToken string) error {
	panic("unimplemented")
}
