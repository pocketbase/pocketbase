package mailer

import (
	"bytes"
	"errors"
	"mime"
	"net/http"
	"os/exec"
	"strings"

	"github.com/pocketbase/pocketbase/tools/hook"
)

var _ Mailer = (*Sendmail)(nil)

// Sendmail implements [mailer.Mailer] interface and defines a mail
// client that sends emails via the "sendmail" *nix command.
//
// This client is usually recommended only for development and testing.
type Sendmail struct {
	onSend *hook.Hook[*SendEvent]
}

// OnSend implements [mailer.SendInterceptor] interface.
func (c *Sendmail) OnSend() *hook.Hook[*SendEvent] {
	if c.onSend == nil {
		c.onSend = &hook.Hook[*SendEvent]{}
	}
	return c.onSend
}

// Send implements [mailer.Mailer] interface.
func (c *Sendmail) Send(m *Message) error {
	if c.onSend != nil {
		return c.onSend.Trigger(&SendEvent{Message: m}, func(e *SendEvent) error {
			return c.send(e.Message)
		})
	}

	return c.send(m)
}

func (c *Sendmail) send(m *Message) error {
	toAddresses := addressesToStrings(m.To, false)

	headers := make(http.Header)
	headers.Set("Subject", mime.QEncoding.Encode("utf-8", m.Subject))
	headers.Set("From", m.From.String())
	headers.Set("Content-Type", "text/html; charset=UTF-8")
	headers.Set("To", strings.Join(toAddresses, ","))

	cmdPath, err := findSendmailPath()
	if err != nil {
		return err
	}

	var buffer bytes.Buffer

	// write
	// ---
	if err := headers.Write(&buffer); err != nil {
		return err
	}
	if _, err := buffer.Write([]byte("\r\n")); err != nil {
		return err
	}
	if m.HTML != "" {
		if _, err := buffer.Write([]byte(m.HTML)); err != nil {
			return err
		}
	} else {
		if _, err := buffer.Write([]byte(m.Text)); err != nil {
			return err
		}
	}
	// ---

	sendmail := exec.Command(cmdPath, strings.Join(toAddresses, ","))
	sendmail.Stdin = &buffer

	return sendmail.Run()
}

func findSendmailPath() (string, error) {
	options := []string{
		"/usr/sbin/sendmail",
		"/usr/bin/sendmail",
		"sendmail",
	}

	for _, option := range options {
		path, err := exec.LookPath(option)
		if err == nil {
			return path, err
		}
	}

	return "", errors.New("failed to locate a sendmail executable path")
}
