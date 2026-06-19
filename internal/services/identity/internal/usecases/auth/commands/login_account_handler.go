package command

import (
	"context"
	"net/http"
	"sync"
	"time"

	apiresponse "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/api-response"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/decorator"
	errorcode "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/error_code"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/jwt"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/security/hash"
	v1dto "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/dtos/v1"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/db/repository"
	"golang.org/x/time/rate"
)

type LoginAccountCommand struct {
	Email    string
	Password string
	ClientIP string
}

type LoginAccountHandler decorator.CommandHandler[LoginAccountCommand, v1dto.LoginAccountResponse]

type loginAccountHandler struct {
	userRepository repository.UserRepository
	jwtService     jwt.JwtService
	hashService    hash.HashService
}

type loginIPRateLimiter struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

var (
	mu                 sync.Mutex
	clientIPs          = make(map[string]*loginIPRateLimiter)
	LOGIN_ATTEMPT_TTL  int
	MAX_LOGIN_ATTEMPTS int
)

func NewLoginAccountHandler(userRepository repository.UserRepository, jwtService jwt.JwtService, hashService hash.HashService) LoginAccountHandler {
	if userRepository == nil {
		panic("nil userRepository")
	}

	if jwtService == nil {
		panic("nil jwtService")
	}

	if hashService == nil {
		panic("nil hashService")
	}

	LOGIN_ATTEMPT_TTL = int(5 * time.Minute)
	MAX_LOGIN_ATTEMPTS = 3

	return decorator.ApplyCommandDecorators(
		&loginAccountHandler{
			userRepository: userRepository,
			jwtService:     jwtService,
			hashService:    hashService,
		},
	)
}

func (h *loginAccountHandler) Handle(ctx context.Context, cmd LoginAccountCommand) (v1dto.LoginAccountResponse, error) {

	if err := h.checkLoginAttempt(cmd.ClientIP); err != nil {
		return v1dto.LoginAccountResponse{}, err
	}

	// validate user exist
	user, err := h.userRepository.GetUserByEmail(ctx, cmd.Email)
	if err != nil {
		h.getLoginAttempt(cmd.ClientIP)

		return v1dto.LoginAccountResponse{}, err
	}

	// verify password
	if err := h.hashService.CompareHashAndPassword(user.UserPassword, cmd.Password); err != nil {
		h.getLoginAttempt(cmd.ClientIP)

		return v1dto.LoginAccountResponse{}, err
	}

	// validate account verification (later)

	// alternate flow:
	// redirect to verify page + token
	// user -> input:email then send email + redirect to verify page
	// user -> input:otp + token

	// generate AT and RT
	accessToken, err := h.jwtService.GenerateAccessToken(jwt.AccessTokenPayload{
		UserUUID:  user.UserUuid.String(),
		UserEmail: user.UserEmail,
	})
	if err != nil {
		return v1dto.LoginAccountResponse{}, err
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(jwt.RefreshTokenPayload{
		UserUUID:  user.UserUuid.String(),
		UserEmail: user.UserEmail,
	})
	if err != nil {
		return v1dto.LoginAccountResponse{}, err
	}

	// store RT to redis
	err = h.jwtService.StoreRefreshToken(refreshToken)
	if err != nil {
		return v1dto.LoginAccountResponse{}, err
	}

	// delete rate limiter
	h.cleanUpClientIP(cmd.ClientIP)

	return v1dto.LoginAccountResponse{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
		ExpiresIn:    accessToken.TTL,
	}, nil
}

func (h *loginAccountHandler) checkLoginAttempt(ip string) error {
	limiter := h.getLoginAttempt(ip)

	if !limiter.Allow() {
		return &apiresponse.ApiError{
			StatusCode: http.StatusTooManyRequests,
			ErrorCode:  errorcode.ErrCodeTooManyRequests,
			Message:    "Too many requests, please try again later.",
		}
	}

	return nil
}

func (h *loginAccountHandler) getLoginAttempt(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clientIPs[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Limit(MAX_LOGIN_ATTEMPTS), MAX_LOGIN_ATTEMPTS)

		newClientIp := &loginIPRateLimiter{
			Limiter:  limiter,
			LastSeen: time.Now(),
		}

		clientIPs[ip] = newClientIp

		return limiter
	}

	client.LastSeen = time.Now()

	return client.Limiter
}

func (h *loginAccountHandler) cleanUpClientIP(ip string) {
	mu.Lock()
	defer mu.Unlock()

	delete(clientIPs, ip)
}
