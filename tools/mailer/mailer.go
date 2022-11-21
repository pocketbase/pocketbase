package mailer

import (
	"io"
	"net/mail"
)

// Message defines a generic email message struct.
type Message struct {
	From        mail.Address
	To          mail.Address
	Bcc         []string
	Cc          []string
	Subject     string
	HTML        string
	Text        string
	Headers     map[string]string
	Attachments map[string]io.Reader
}

// Mailer defines a base mail client interface.
type Mailer interface {
	// Send sends an email with the provided Message.
	Send(message *Message) error
}
