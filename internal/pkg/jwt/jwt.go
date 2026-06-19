package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/security/encrypt"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/config"
)

type JwtService interface {
	GenerateAccessToken(payload AccessTokenPayload) (AccessToken, error)
	GenerateRefreshToken(payload RefreshTokenPayload) (RefreshToken, error)
	StoreRefreshToken(refreshToken RefreshToken) error
}

type AccessToken struct {
	Token string
	TTL   int
}

type AccessTokenPayload struct {
	UserUUID  string
	UserEmail string
}

type RefreshTokenPayload struct {
	UserUUID  string
	UserEmail string
}

type RefreshToken struct {
	Token     string
	UserUUID  string
	UserEmail string
	ExpireAt  time.Time
	Revoked   bool
}

var (
	JWT_SECRET            string
	JWT_ENCRYPTION_KEY    string
	JWT_ACCESS_TOKEN_TTL  time.Duration
	JWT_REFRESH_TOKEN_TTL time.Duration
	JWT_ISSUER            string

	CACHE_REFRESH_TOKEN_KEY_PREFIX string = "REFRESH_TOKEN:"
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

func (jv *jwtService) GenerateAccessToken(payload AccessTokenPayload) (AccessToken, error) {

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

func (jv *jwtService) GenerateRefreshToken(payload RefreshTokenPayload) (RefreshToken, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return RefreshToken{}, err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	return RefreshToken{
		Token:     token,
		UserUUID:  payload.UserUUID,
		UserEmail: payload.UserEmail,
		ExpireAt:  time.Now().Add(JWT_REFRESH_TOKEN_TTL),
		Revoked:   false,
	}, nil
}

func (jv *jwtService) StoreRefreshToken(refreshToken RefreshToken) error {
	panic("unimplemented")
}
