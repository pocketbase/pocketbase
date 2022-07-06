package tests

import (
	"io"
	"net/mail"

	"github.com/pocketbase/pocketbase/tools/mailer"
)

var _ mailer.Mailer = (*TestMailer)(nil)

// TestMailer is a mock `mailer.Mailer` implementation.
type TestMailer struct {
	TotalSend    int
	LastHtmlBody string
}

// Reset clears any previously test collected data.
func (m *TestMailer) Reset() {
	m.LastHtmlBody = ""
	m.TotalSend = 0
}

// Send implements `mailer.Mailer` interface.
func (m *TestMailer) Send(fromEmail mail.Address, toEmail mail.Address, subject string, html string, attachments map[string]io.Reader) error {
	m.LastHtmlBody = html
	m.TotalSend++
	return nil
}
