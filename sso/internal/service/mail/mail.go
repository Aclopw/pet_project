package mail

import (
	"fmt"
	"log/slog"
	"sso/internal/config"

	gomail "gopkg.in/mail.v2"
)

type MailService struct {
	log        *slog.Logger
	mailConfig *config.MailServer
}

func New(log *slog.Logger, mailConfig *config.MailServer) *MailService {
	return &MailService{
		log:        log,
		mailConfig: mailConfig,
	}
}

func (m *MailService) Send(email, activationLink string) error {
	const op = "mail.MailService.Send"

	m.log.Info("Sending verification email", slog.String("email", email), slog.String("activationLink", activationLink))

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.mailConfig.Username)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Please verify your email address")
	msg.SetBody("text/plain", "Click the link to verify your email: "+activationLink)

	dialer := gomail.NewDialer(m.mailConfig.Host, m.mailConfig.Port, m.mailConfig.Username, m.mailConfig.Password)

	if err := dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	m.log.Info("Verification email sent successfully", slog.String("email", email))

	return nil
}
