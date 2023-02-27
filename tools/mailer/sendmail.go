package mailer

import (
	"bytes"
	"errors"
	"mime"
	"net/http"
	"os/exec"
	"strings"
)

var _ Mailer = (*Sendmail)(nil)

// Sendmail implements [mailer.Mailer] interface and defines a mail
// client that sends emails via the "sendmail" *nix command.
//
// This client is usually recommended only for development and testing.
type Sendmail struct {
}

// Send implements `mailer.Mailer` interface.
func (c *Sendmail) Send(m *Message) error {
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

	return "", errors.New("Failed to locate a sendmail executable path.")
}
