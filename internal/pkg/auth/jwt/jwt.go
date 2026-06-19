package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/cache"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/security/encrypt"
)

type JwtService interface {
	GenerateAccessToken(payload AccessTokenPayload) (AccessToken, error)
	GenerateRefreshToken(payload RefreshTokenPayload) (RefreshToken, error)
	StoreRefreshToken(refreshToken RefreshToken) error
	ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error)
	DecryptAccessTokenPayload(tokenString string) (*AccessTokenPayload, error)
	ValidateRefreshToken(token string) (RefreshToken, error)
	RevokeRefreshToken(token string) error
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

	CACHE_REFRESH_TOKEN_KEY_PREFIX string = "REFRESH_TOKEN"
)

type jwtService struct {
	redisService cache.CacheService
}

// !!! DO NOT IMPORT "internal/services/identity/config" HERE !!!
// This is a shared package (root module). Importing a specific service (sub-module)
// creates a circular dependency that breaks go mod tidy!
func NewJwtService(redisService cache.CacheService, secret, encryptionKey, issuer string, accessTTL, refreshTTL time.Duration) JwtService {

	if redisService == nil {
		panic("Redis service is nil")
	}

	JWT_SECRET = secret
	JWT_ENCRYPTION_KEY = encryptionKey
	JWT_ACCESS_TOKEN_TTL = accessTTL
	JWT_REFRESH_TOKEN_TTL = refreshTTL
	JWT_ISSUER = issuer

	return &jwtService{
		redisService: redisService,
	}
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
	cacheKey := fmt.Sprintf("%s:%s", CACHE_REFRESH_TOKEN_KEY_PREFIX, refreshToken.Token)

	err := jv.redisService.Set(cacheKey, refreshToken, JWT_REFRESH_TOKEN_TTL)

	if err != nil {
		return err
	}

	return nil
}

func (jv *jwtService) DecryptAccessTokenPayload(tokenString string) (*AccessTokenPayload, error) {
	panic("unimplemented")
}

func (jv *jwtService) ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	panic("unimplemented")
}

func (jv *jwtService) RevokeRefreshToken(token string) error {
	panic("unimplemented")
}

func (jv *jwtService) ValidateRefreshToken(token string) (RefreshToken, error) {
	panic("unimplemented")
}
