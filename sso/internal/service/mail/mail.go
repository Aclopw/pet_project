package mail

import (
	"log/slog"
	"net/smtp"
	"sso/internal/config"
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

	auth := smtp.PlainAuth(
		"",
		m.mailConfig.Username,
		m.mailConfig.Password,
		m.mailConfig.Host,
	)

	to := []string{email}

	message := []byte(
		"Subject: Verify your email\n" + "activation link: " + activationLink)

	err := smtp.SendMail(m.mailConfig.Host, auth, m.mailConfig.Username, to, message)

	if err != nil {
		m.log.Error("%s: %w", op, err)

		return err
	}

	m.log.Info("Verification email sent successfully", slog.String("email", email))

	return nil
}
