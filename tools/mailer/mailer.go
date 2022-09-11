package mailer

import (
	"io"
	"net/mail"
)

// Mailer defines a base mail client interface.
type Mailer interface {
	// Send sends an email with HTML body to the specified recipient.
	Send(
		fromEmail mail.Address,
		toEmail mail.Address,
		subject string,
		htmlContent string,
		attachments map[string]io.Reader,
	) error
}
