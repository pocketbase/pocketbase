package mailer

import (
	"fmt"
	"io"
	"net/mail"
	"net/smtp"
	"regexp"
	"strings"

	"github.com/domodwyer/mailyak/v3"
	"github.com/microcosm-cc/bluemonday"
)

var _ Mailer = (*SmtpClient)(nil)

// regex to select all tabs
var tabsRegex = regexp.MustCompile(`\t+`)

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
func (m *SmtpClient) Send(
	fromEmail mail.Address,
	toEmail mail.Address,
	subject string,
	htmlBody string,
	attachments map[string]io.Reader,
) error {
	smtpAuth := smtp.PlainAuth("", m.username, m.password, m.host)

	// create mail instance
	var yak *mailyak.MailYak
	if m.tls {
		var tlsErr error
		yak, tlsErr = mailyak.NewWithTLS(fmt.Sprintf("%s:%d", m.host, m.port), smtpAuth, nil)
		if tlsErr != nil {
			return tlsErr
		}
	} else {
		yak = mailyak.New(fmt.Sprintf("%s:%d", m.host, m.port), smtpAuth)
	}

	if fromEmail.Name != "" {
		yak.FromName(fromEmail.Name)
	}
	yak.From(fromEmail.Address)
	yak.To(toEmail.Address)
	yak.Subject(subject)
	yak.HTML().Set(htmlBody)

	// set also plain text content
	policy := bluemonday.StrictPolicy() // strips all tags
	yak.Plain().Set(strings.TrimSpace(tabsRegex.ReplaceAllString(policy.Sanitize(htmlBody), "")))

	for name, data := range attachments {
		yak.Attach(name, data)
	}

	return yak.Send()
}
