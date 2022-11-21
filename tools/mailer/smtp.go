package mailer

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/domodwyer/mailyak/v3"
	"github.com/pocketbase/pocketbase/tools/security"
)

var _ Mailer = (*SmtpClient)(nil)

// NewSmtpClient creates new `SmtpClient` with the provided configuration.
func NewSmtpClient(
	host string,
	port int,
	username string,
	password string,
	tls bool,
) *SmtpClient {
	return &SmtpClient{
		host:     host,
		port:     port,
		username: username,
		password: password,
		tls:      tls,
	}
}

// SmtpClient defines a SMTP mail client structure that implements
// `mailer.Mailer` interface.
type SmtpClient struct {
	host     string
	port     int
	username string
	password string
	tls      bool
}

// Send implements `mailer.Mailer` interface.
func (c *SmtpClient) Send(m *Message) error {
	var smtpAuth smtp.Auth
	if c.username != "" || c.password != "" {
		smtpAuth = smtp.PlainAuth("", c.username, c.password, c.host)
	}

	// create mail instance
	var yak *mailyak.MailYak
	if c.tls {
		var tlsErr error
		yak, tlsErr = mailyak.NewWithTLS(fmt.Sprintf("%s:%d", c.host, c.port), smtpAuth, nil)
		if tlsErr != nil {
			return tlsErr
		}
	} else {
		yak = mailyak.New(fmt.Sprintf("%s:%d", c.host, c.port), smtpAuth)
	}

	if m.From.Name != "" {
		yak.FromName(m.From.Name)
	}
	yak.From(m.From.Address)
	yak.To(m.To.Address)
	yak.Subject(m.Subject)
	yak.HTML().Set(m.HTML)

	if m.Text == "" {
		// try to generate a plain text version of the HTML
		if plain, err := html2Text(m.HTML); err == nil {
			yak.Plain().Set(plain)
		}
	} else {
		yak.Plain().Set(m.Text)
	}

	if len(m.Bcc) > 0 {
		yak.Bcc(m.Bcc...)
	}

	if len(m.Cc) > 0 {
		yak.Cc(m.Cc...)
	}

	// add attachements (if any)
	for name, data := range m.Attachments {
		yak.Attach(name, data)
	}

	// add custom headers (if any)
	var hasMessageId bool
	for k, v := range m.Headers {
		if strings.EqualFold(k, "Message-ID") {
			hasMessageId = true
		}
		yak.AddHeader(k, v)
	}
	if !hasMessageId {
		// add a default message id if missing
		fromParts := strings.Split(m.From.Address, "@")
		if len(fromParts) == 2 {
			yak.AddHeader("Message-ID", fmt.Sprintf("<%s@%s>",
				security.PseudorandomString(15),
				fromParts[1],
			))
		}
	}

	return yak.Send()
}
