package command

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	apiresponse "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/api-response"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/cache"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/decorator"
	errorcode "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/error_code"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/queue/rabbitmq"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/config"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/external/mail"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/utils"
	"golang.org/x/time/rate"
)

type ForgotPasswordCommand struct {
	Email string
}

type ForgotPasswordHandler decorator.CommandHandler[ForgotPasswordCommand, bool]

type sendEmailAttempt struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	sendEmailMutex  sync.Mutex
	sendEmailEmails = make(map[string]*sendEmailAttempt)
)

const (
	SENT_EMAIL_ATTEMPT_TTL  = 3 * time.Minute
	MAX_SEND_EMAIL_ATTEMPTS = 3
)

type forgotPasswordHandler struct {
	mq          rabbitmq.MessageQueueService
	cache       cache.CacheService
	mailService mail.MailService
}

func NewForgotPasswordHandler(mq rabbitmq.MessageQueueService, cache cache.CacheService, mailService mail.MailService) ForgotPasswordHandler {

	if mq == nil {
		panic("nil message queue service")
	}

	if cache == nil {
		panic("nil cache service")
	}

	if mailService == nil {
		panic("nil mail service")
	}

	return decorator.ApplyCommandDecorators(
		&forgotPasswordHandler{
			mq:          mq,
			cache:       cache,
			mailService: mailService,
		},
	)
}

func (h *forgotPasswordHandler) Handle(ctx context.Context, cmd ForgotPasswordCommand) (bool, error) {

	// create rate limiter
	if err := h.checkEmailAttempt(cmd.Email); err != nil {
		return false, err
	}

	// validate email exist

	// create reset token
	token, err := utils.GenerateRandomString(16)
	if err != nil {
		return false, err
	}

	// store reset token to redis
	resetPasswordKey := "RESET_PASSWORD:" + token
	if err := h.cache.Set(resetPasswordKey, "userid", 5*time.Minute); err != nil {
		return false, err
	}

	// send email with reset token
	resetLink := fmt.Sprintf("https://localhost:3000/auth/verify-reset-password?token=%s", token)

	mailContent := &mail.Email{
		From: mail.Address{
			Name:  "YGZ",
			Email: config.AppConfig.Mail.Sender,
		},
		To: []mail.Address{
			{
				Email: cmd.Email,
			},
		},
		Subject: "YGZ RESET PASSWORD",
		Text: fmt.Sprintf("Hi %s, \n\n You requested to reset your password. Please click the link below to reset it:\n%s\n\n The link will expire in 1 hour. \n\n Best regard, \nCode With Tuan Team",
			cmd.Email,
			resetLink),
	}

	if err := h.mailService.SendMail(ctx, mailContent); err != nil {
		return false, err
	}

	// if err := h.mq.Publish("reset_password_queue", mailContent); err != nil {
	// 	return false, err
	// }

	return true, nil
}

func (h *forgotPasswordHandler) checkEmailAttempt(email string) error {
	attemptEmail := h.getSendEmailAttempt(email)

	if !attemptEmail.limiter.Allow() {
		return &apiresponse.ApiError{
			StatusCode: http.StatusTooManyRequests,
			ErrorCode:  errorcode.ErrCodeTooManyRequests,
			Message:    "Too many attempts, please try again later.",
		}
	}

	return nil
}

func (h *forgotPasswordHandler) getSendEmailAttempt(email string) *sendEmailAttempt {
	sendEmailMutex.Lock()
	defer sendEmailMutex.Unlock()

	attempt, exists := sendEmailEmails[email]
	if !exists {
		limiter := rate.NewLimiter(rate.Limit(float32(MAX_SEND_EMAIL_ATTEMPTS)/float32(SENT_EMAIL_ATTEMPT_TTL.Seconds())), MAX_SEND_EMAIL_ATTEMPTS)
		newAttempt := &sendEmailAttempt{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		sendEmailEmails[email] = newAttempt
		return newAttempt
	}

	attempt.lastSeen = time.Now()
	return attempt
}
