package mail

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type MailProvider string

const (
	GoogleMailProvider MailProvider = "google"
)

type Email struct {
	From     Address   `json:"from"`
	To       []Address `json:"to"`
	Subject  string    `json:"subject"`
	Text     string    `json:"text"`
	Category string    `json:"category"`
}

type Address struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type MailService interface {
	SendMail(ctx context.Context, email *Email) error
}

type MailConfig struct {
	MaxRetries uint
	Timeout    time.Duration
	Logger     *zerolog.Logger
}

type mailService struct {
	config   MailConfig
	provider MailService
	logger   *zerolog.Logger
}

func NewMailService(providerFactory MailFactory, logger *zerolog.Logger) MailService {

	config := &MailConfig{
		MaxRetries: 3,
		Timeout:    10 * time.Second,
		Logger:     logger,
	}

	provider := providerFactory.CreateProvider(logger)

	return &mailService{
		config:   *config,
		logger:   logger,
		provider: provider,
	}
}

func (ms *mailService) SendMail(ctx context.Context, email *Email) error {

	traceId := ctx.Value("trace_id").(string)

	start := time.Now()

	var lastError error

	for attempt := 1; attempt <= int(ms.config.MaxRetries); attempt++ {
		startAttempt := time.Now()

		err := ms.provider.SendMail(ctx, email)
		if err == nil {
			ms.logger.Info().Str("trace_id", traceId).
				Dur("duration", time.Since(startAttempt)).
				Str("operation", "send_mail").
				Interface("to", email.To).
				Str("subject", email.Subject).
				Str("category", email.Category).
				Msg("Email send successfully")

			return nil
		}

		lastError = err
		ms.logger.Warn().Str("trace_id", traceId).
			Dur("duration", time.Since(startAttempt)).
			Str("operation", "send_mail").
			Int("attempt", attempt).
			Err(err).
			Msg("Failed to send email, retrying")

		time.Sleep(time.Duration(attempt) * time.Second)
	}

	ms.logger.Error().Str("trace_id", traceId).
		Dur("duration", time.Since(start)).
		Str("operation", "send_mail").
		Int("attempt", int(ms.config.MaxRetries)).
		Err(lastError).
		Msg("Failed to send email after all retries")

	return nil
}
