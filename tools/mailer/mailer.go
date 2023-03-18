package mailer

import (
	"io"
	"net/mail"
)

// Message defines a generic email message struct.
type Message struct {
	From        mail.Address
	To          []mail.Address
	Bcc         []mail.Address
	Cc          []mail.Address
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

// addressesToStrings converts the provided address to a list of serialized RFC 5322 strings.
//
// To export only the email part of mail.Address, you can set withName to false.
func addressesToStrings(addresses []mail.Address, withName bool) []string {
	result := make([]string, len(addresses))

	for i, addr := range addresses {
		if withName && addr.Name != "" {
			result[i] = addr.String()
		} else {
			// keep only the email part to avoid wrapping in angle-brackets
			result[i] = addr.Address
		}
	}

	return result
}
