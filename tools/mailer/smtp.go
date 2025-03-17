package mailer

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/domodwyer/mailyak/v3"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/security"
)

var _ Mailer = (*SMTPClient)(nil)

const (
	SMTPAuthPlain = "PLAIN"
	SMTPAuthLogin = "LOGIN"
)

// SMTPClient defines a SMTP mail client structure that implements
// `mailer.Mailer` interface.
type SMTPClient struct {
	onSend *hook.Hook[*SendEvent]

	TLS      bool
	Port     int
	Host     string
	Username string
	Password string

	// SMTP auth method to use
	// (if not explicitly set, defaults to "PLAIN")
	AuthMethod string

	// LocalName is optional domain name used for the EHLO/HELO exchange
	// (if not explicitly set, defaults to "localhost").
	//
	// This is required only by some SMTP servers, such as Gmail SMTP-relay.
	LocalName string
}

// OnSend implements [mailer.SendInterceptor] interface.
func (c *SMTPClient) OnSend() *hook.Hook[*SendEvent] {
	if c.onSend == nil {
		c.onSend = &hook.Hook[*SendEvent]{}
	}
	return c.onSend
}

// Send implements [mailer.Mailer] interface.
func (c *SMTPClient) Send(m *Message) error {
	if c.onSend != nil {
		return c.onSend.Trigger(&SendEvent{Message: m}, func(e *SendEvent) error {
			return c.send(e.Message)
		})
	}

	return c.send(m)
}

func (c *SMTPClient) send(m *Message) error {
	var smtpAuth smtp.Auth
	if c.Username != "" || c.Password != "" {
		switch c.AuthMethod {
		case SMTPAuthLogin:
			smtpAuth = &smtpLoginAuth{c.Username, c.Password}
		default:
			smtpAuth = smtp.PlainAuth("", c.Username, c.Password, c.Host)
		}
	}

	// create mail instance
	var yak *mailyak.MailYak
	if c.TLS {
		var tlsErr error
		yak, tlsErr = mailyak.NewWithTLS(fmt.Sprintf("%s:%d", c.Host, c.Port), smtpAuth, nil)
		if tlsErr != nil {
			return tlsErr
		}
	} else {
		yak = mailyak.New(fmt.Sprintf("%s:%d", c.Host, c.Port), smtpAuth)
	}

	if c.LocalName != "" {
		yak.LocalName(c.LocalName)
	}

	if m.From.Name != "" {
		yak.FromName(m.From.Name)
	}
	yak.From(m.From.Address)
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

	if len(m.To) > 0 {
		yak.To(addressesToStrings(m.To, true)...)
	}

	if len(m.Bcc) > 0 {
		yak.Bcc(addressesToStrings(m.Bcc, true)...)
	}

	if len(m.Cc) > 0 {
		yak.Cc(addressesToStrings(m.Cc, true)...)
	}

	// add regular attachements (if any)
	for name, data := range m.Attachments {
		r, mime, err := detectReaderMimeType(data)
		if err != nil {
			return err
		}
		yak.AttachWithMimeType(name, r, mime)
	}

	// add inline attachments (if any)
	for name, data := range m.InlineAttachments {
		r, mime, err := detectReaderMimeType(data)
		if err != nil {
			return err
		}
		yak.AttachInlineWithMimeType(name, r, mime)
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

// -------------------------------------------------------------------
// AUTH LOGIN
// -------------------------------------------------------------------

var _ smtp.Auth = (*smtpLoginAuth)(nil)

// smtpLoginAuth defines an AUTH that implements the LOGIN authentication mechanism.
//
// AUTH LOGIN is obsolete[1] but some mail services like outlook requires it [2].
//
// NB!
// It will only send the credentials if the connection is using TLS or is connected to localhost.
// Otherwise authentication will fail with an error, without sending the credentials.
//
// [1]: https://github.com/golang/go/issues/40817
// [2]: https://support.microsoft.com/en-us/office/outlook-com-no-longer-supports-auth-plain-authentication-07f7d5e9-1697-465f-84d2-4513d4ff0145?ui=en-us&rs=en-us&ad=us
type smtpLoginAuth struct {
	username, password string
}

// Start initializes an authentication with the server.
//
// It is part of the [smtp.Auth] interface.
func (a *smtpLoginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	// Must have TLS, or else localhost server.
	// Note: If TLS is not true, then we can't trust ANYTHING in ServerInfo.
	// In particular, it doesn't matter if the server advertises LOGIN auth.
	// That might just be the attacker saying
	// "it's ok, you can trust me with your password."
	if !server.TLS && !isLocalhost(server.Name) {
		return "", nil, errors.New("unencrypted connection")
	}

	return "LOGIN", nil, nil
}

// Next "continues" the auth process by feeding the server with the requested data.
//
// It is part of the [smtp.Auth] interface.
func (a *smtpLoginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch strings.ToLower(string(fromServer)) {
		case "username:":
			return []byte(a.username), nil
		case "password:":
			return []byte(a.password), nil
		}
	}

	return nil, nil
}

func isLocalhost(name string) bool {
	return name == "localhost" || name == "127.0.0.1" || name == "::1"
}
