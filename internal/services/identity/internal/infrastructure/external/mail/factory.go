package mail

import (
	"fmt"

	"github.com/rs/zerolog"
)

type MailFactory interface {
	CreateProvider(logger *zerolog.Logger) MailService
}

type GoogleMailFactory struct{}

func (f *GoogleMailFactory) CreateProvider(logger *zerolog.Logger) MailService {
	return NewGoogleMail(logger)
}

func NewMailFactory(provider MailProvider) (MailFactory, error) {

	switch provider {
	case GoogleMailProvider:
		return &GoogleMailFactory{}, nil
	default:
		return nil, fmt.Errorf("mail provider %s not found", provider)
	}
}
