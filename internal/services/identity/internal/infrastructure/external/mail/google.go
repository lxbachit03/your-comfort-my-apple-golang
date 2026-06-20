package mail

import (
	"context"
	"log"

	"net/smtp"

	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/config"
	"github.com/rs/zerolog"
)

type GoogleMail struct {
	config *GoogleMailConfig
	logger *zerolog.Logger
}

type GoogleMailConfig struct {
	MailSender string
	NameSender string
}

func NewGoogleMail(logger *zerolog.Logger) MailService {
	return &GoogleMail{
		config: &GoogleMailConfig{
			MailSender: config.AppConfig.Mail.Sender,
			NameSender: "ygz",
		},
		logger: logger,
	}

}

func (g GoogleMail) SendMail(ctx context.Context, email *Email) error {

	from := config.AppConfig.Mail.Sender
	password := config.AppConfig.Mail.Provider.Google.AppPassword

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	var toEmails []string
	var toHeader string
	for i, addr := range email.To {
		toEmails = append(toEmails, addr.Email)
		if i > 0 {
			toHeader += ", "
		}
		toHeader += addr.Email
	}

	msg := "From: " + from + "\r\n" +
		"To: " + toHeader + "\r\n" +
		"Subject: " + email.Subject + "\r\n\r\n" +
		email.Text

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, toEmails, []byte(msg))
	if err != nil {
		log.Printf("google send mail failed: %v", err)
		return err
	}

	log.Printf("google send mail success to: %v", toEmails)
	return nil
}
