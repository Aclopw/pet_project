package mail

import "log/slog"

type MailService struct {
	log *slog.Logger
}

func New(log *slog.Logger) *MailService {
	return &MailService{
		log: log,
	}
}

func (m *MailService) Send(email, activationLink string) error {
	m.log.Info("sending verification email", slog.String("email", email), slog.String("activation_link", activationLink))

	// Here you would implement the actual email sending logic.
	// For example, using an SMTP client or a third-party service.

	// Simulating email sending for demonstration purposes.
	m.log.Info("verification email sent successfully", slog.String("email", email))

	return nil
}
